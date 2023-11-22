package shgui

import (
	"embed"
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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

	if dev.IsMobile() {
		// Mobile platform
		spacePosition = layout.NewVBoxLayout()
		dialogSize = fyne.NewSize(300, 600)
	} else {
		// PC platform
		// spaceposition = layout.NewHBoxLayout()
		spacePosition = layout.NewVBoxLayout()
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
	// Create button for key file open
	var keyButton *widget.Button
	var keyURI fyne.URIReadCloser
	keyButton = widget.NewButton("Key file",
		func() {
			win := a.NewWindow("Open key file")
			win.Resize(dialogSize)
			win.SetFixedSize(true)
			win.Show()
			keyButton.Disable()
			doFileOpen := func(f fyne.URIReadCloser, err error) {
				if err != nil {
					log.Printf("Error opening file: %s", err)
				} else {
					keyURI = f
				}
				keyButton.Enable()
				win.Close()
				// TODO: Delete
				if keyURI != nil {
					fmt.Println(keyURI.URI())
				}
			}
			dlg := dialog.NewFileOpen(doFileOpen, win)
			dlg.Resize(dialogSize)
			dlg.Show()
		})
	// Create button for source file open
	var srcButton *widget.Button
	var srcURI fyne.URIReadCloser
	srcButton = widget.NewButton("Source file",
		func() {
			win := a.NewWindow("Open source file")
			win.Resize(dialogSize)
			win.SetFixedSize(true)
			win.Show()
			srcButton.Disable()
			doFileOpen := func(f fyne.URIReadCloser, err error) {
				if err != nil {
					log.Printf("Error opening file: %s", err)
				} else {
					srcURI = f
				}
				srcButton.Enable()
				win.Close()
				// TODO: Delete
				if srcURI != nil {
					fmt.Println(srcURI.URI())
				}
			}
			dlg := dialog.NewFileOpen(doFileOpen, win)
			dlg.Resize(dialogSize)
			dlg.Show()
		})
	// Create button for output file
	var outButton *widget.Button
	var outURI fyne.URIWriteCloser
	outButton = widget.NewButton("Output file",
		func() {
			win := a.NewWindow("Save file")
			win.Resize(dialogSize)
			win.SetFixedSize(true)
			win.Show()
			outButton.Disable()
			doFileSave := func(f fyne.URIWriteCloser, err error) {
				if err != nil {
					log.Printf("Error opening file: %s", err)
				} else {
					outURI = f
				}
				outButton.Enable()
				win.Close()
				// TODO: Delete
				if outURI != nil {
					fmt.Println(outURI.URI())
				}
			}
			dlg := dialog.NewFileSave(doFileSave, win)
			dlg.Resize(dialogSize)
			dlg.Show()
		})
	// Make GUI
	appcontainer := fyne.NewContainer()
	appcontainer = fyne.NewContainerWithLayout(spacePosition, appimage)
	subcontainer := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		keyButton,
		srcButton,
		outButton)
	appcontainer.AddObject(subcontainer)

	w.SetContent(appcontainer)
	w.Show()
	a.Run()
}
