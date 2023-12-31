package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shdtd/shcrypt/lib/shgui"
	"github.com/shdtd/shcrypt/lib/shresource"
)

func main() {
	// Setting commands and program arguments
	// Setting "encrypt" command
	encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
	fileForEncrypt := encryptCmd.String("file", "", "file")
	passphraseFileOnEncrypt := encryptCmd.String("key", "", "key")
	outForEncrypt := encryptCmd.String("out", "", "out")
	// Setting "decrypt" command
	decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
	fileForDecrypt := decryptCmd.String("file", "", "file")
	passphraseFileOnDecrypt := decryptCmd.String("key", "", "key")
	outForDecrypt := decryptCmd.String("out", "", "out")

	// If the program is run without arguments, the GUI will be displayed
	if len(os.Args) < 2 {
		shgui.Run()
		os.Exit(0)
	}
	// Loop through all arguments
	switch os.Args[1] {
	// Handling the "encrypt" command
	case "encrypt":
		// Parsing the encryption command arguments
		encryptCmd.Parse(os.Args[2:])
		// Encryption
		res, err := shresource.NewShResource(*fileForEncrypt,
			*outForEncrypt,
			*passphraseFileOnEncrypt)
		if err != nil {
			fmt.Println("Get resources returned an error:", err)
			os.Exit(0)
		}

		res.Type = "encrypt"
		err = res.Encrypt()
		if err != nil {
			fmt.Println("Encrypt error:", err)
			os.Exit(0)
		}
		// Write encrypted data to file
		res.FileSafe()
	// Handling the "decrypt" command
	case "decrypt":
		// Parsing the decryption command arguments
		decryptCmd.Parse(os.Args[2:])
		// Decryption
		res, err := shresource.NewShResource(*fileForDecrypt,
			*outForDecrypt,
			*passphraseFileOnDecrypt)
		if err != nil {
			fmt.Println("Get resources returned an error:", err)
			os.Exit(0)
		}

		res.Type = "decrypt"
		err = res.Decrypt()
		if err != nil {
			fmt.Println("Decrypt error:", err)
			os.Exit(0)
		}
		// Write decrypted data to file
		res.FileSafe()
	// By default, displays help information and exits the program
	default:
		doHelp()
	}
}

// The function displays help information and exits the program
func doHelp() {
	fmt.Println("Expected 'encrypt' or 'decrypt' command")
	fmt.Println("Example:")
	fmt.Printf("For encrypt: %s encrypt --file \"unencrypted_file.xxx\" "+
		"--key \"passphrase_file\" --out \"output_file\"\n",
		os.Args[0])
	fmt.Printf("For decrypt: %s decrypt --file \"encrypted_file.xxx.shc\" "+
		"--key \"passphrase_file\" --out \"output_file\"\n",
		os.Args[0])
	fmt.Println("To display the GUI: Run the program without any arguments.")
	os.Exit(0)
}
