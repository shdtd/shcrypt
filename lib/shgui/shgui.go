package shgui

import (
	"embed"
	"fmt"
	"image/color"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/shdtd/shcrypt/lib/shresource"
)

//go:embed all:res
var resources embed.FS

func Run() {
	// Application, main window and device.
	a := app.New()
	w := a.NewWindow("ShCrypt")
	dev := fyne.CurrentDevice()

	// Read logo image.
	logo, err := resources.ReadFile("res/icons/logo.png")
	if err != nil {
		log.Fatalf("Error reading resource for logo: %s", err)
	}

	// Create fyne.Resource and set application and window icons.
	app_icon := fyne.NewStaticResource("Application icon", logo)
	a.SetIcon(app_icon)
	w.SetIcon(app_icon)

	// Variables dependent on screen position
	var spacePosition fyne.Layout
	var dialogSize fyne.Size

	// Determining the platform
	if dev.IsMobile() {
		// Mobile platform
		spacePosition = layout.NewVBoxLayout()
		dialogSize = fyne.NewSize(300, 600)
	} else {
		// PC platform
		spacePosition = layout.NewHBoxLayout()
		// spacePosition = layout.NewVBoxLayout()
		dialogSize = fyne.NewSize(600, 400)
	}

	// Create application image
	f, err := resources.Open("res/icons/image.png")
	if err != nil {
		log.Fatalf("Error open resource: %s", err)
	}
	defer f.Close()
	appimage := canvas.NewImageFromReader(f, "Application image")
	appimage.FillMode = canvas.ImageFillContain
	appimage.SetMinSize(fyne.NewSize(300, 300))
	appimage.Resize(fyne.NewSize(300, 300))

	// Create label for test
	applabel := canvas.NewText("Encryption/Decryption", color.CMYK{140, 70, 35, 0})
	applabel.TextSize = 24
	applabel.Alignment = fyne.TextAlignCenter
	applabel.TextStyle = fyne.TextStyle{Bold: true}

	// Create object of interface File URI
	files := &FilesURI{}

	// Create button for key file open
	var keyButton *widget.Button
	keyButton = widget.NewButton("Key file",
		func() {
			fileDialog(a, keyButton, dialogSize, files.SetKeyURI, false, "Open key file")
		})

	// Create button for source file open
	var srcButton *widget.Button
	srcButton = widget.NewButton("Source file",
		func() {
			fileDialog(a, srcButton, dialogSize, files.SetSrcURI, false, "Open source file")
		})

	// Create button for output file
	var outButton *widget.Button
	outButton = widget.NewButton("Output file",
		func() {
			fileDialog(a, outButton, dialogSize, files.SetOutURI, true, "Save file")
		})

	// Create button for encryption
	encryptButton := widget.NewButton("Encrypt",
		func() {
			go encrypt(files.SrcURI, files.OutURI, files.KeyURI)
		})

	// Create button for decryption
	decryptButton := widget.NewButton("Decrypt",
		func() {
			go decrypt(files.SrcURI, files.OutURI, files.KeyURI)
		})

	// Create button for exit
	exitButton := widget.NewButton("Exit",
		func() {
			a.Quit()
		})

	// Create space for buttons
	space := container.New(
		layout.NewVBoxLayout(),
		keyButton,
		srcButton,
		outButton,
		encryptButton,
		decryptButton,
		exitButton,
	)
	// Make GUI
	appcontainer := container.New(spacePosition, appimage)
	subcontainer := container.New(layout.NewVBoxLayout(),
		keyButton,
		srcButton,
		outButton)
	_ = subcontainer
	appcontainer.Add(space)
	w.SetContent(appcontainer)
	w.ShowAndRun()
}

func fileDialog(a fyne.App,
	activeButton *widget.Button,
	dialogSize fyne.Size,
	setter func(string),
	isWrite bool,
	title string,
) {
	// Create window for the dialog
	dlgWin := a.NewWindow(title)
	dlgWin.SetContent(container.NewVBox())
	dlgWin.Resize(dialogSize)
	dlgWin.SetFixedSize(true)
	dlgWin.Show()
	activeButton.Disable()

	dlgWin.SetOnClosed(func() {
		activeButton.Enable()
	})

	doFileSave := func(f fyne.URIWriteCloser, err error) {
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		} else {
			setter(f.URI().Path())
		}
		activeButton.Enable()
		dlgWin.Close()
	}

	doFileOpen := func(f fyne.URIReadCloser, err error) {
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		} else {
			setter(f.URI().Path())
		}
		activeButton.Enable()
		dlgWin.Close()
	}

	var dlg *dialog.FileDialog

	switch isWrite {
	default:
		fmt.Printf("Unexpected type: %T\n", isWrite)
	case true:
		dlg = dialog.NewFileSave(doFileSave, dlgWin)
	case false:
		dlg = dialog.NewFileOpen(doFileOpen, dlgWin)
	}

	dlg.Resize(dialogSize)
	dlg.Show()
}

func decrypt(srcURI, outURI, keyURI string) {
	res, err := shresource.NewShResource(srcURI, outURI, keyURI)
	if err != nil {
		fmt.Println("Get resources returned an error:", err)
		os.Exit(1)
	}

	res.Type = "decrypt"
	err = res.Decrypt()
	if err != nil {
		fmt.Println("Dencrypt error:", err)
		os.Exit(1)
	}
	// Write decrypted data to file
	res.FileSafe()
}

func encrypt(srcURI, outURI, keyURI string) {
	res, err := shresource.NewShResource(srcURI, outURI, keyURI)
	if err != nil {
		fmt.Println("Get resources returned an error:", err)
		os.Exit(1)
	}

	res.Type = "encrypt"
	err = res.Encrypt()
	if err != nil {
		fmt.Println("Encrypt error:", err)
		os.Exit(1)
	}
	// Write encrypted data to file
	res.FileSafe()
}
