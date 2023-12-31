package shresource

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type shResource struct {
	source     string
	out        string
	passphrase string
	padding    int
	Type       string
	Block      cipher.Block
	Data       []byte
}

func (r shResource) GetSource() string {
	return r.source
}

func (r shResource) GetOut() string {
	return r.out
}

func NewShResource(source string, out string, passphrase string) (*shResource, error) {
	res := &shResource{}
	res.source = source
	res.out = out
	res.passphrase = passphrase
	err := res.readFiles()
	if err != nil {
		log.Fatalf("Reading files returned an error: %s", err)
		return &shResource{}, err
	}

	return res, err
}

func (res *shResource) readFiles() error {
	keyFile, err := filepath.Abs(res.passphrase)
	if err != nil {
		log.Fatalf("Absolute path for key not found: %s", err)
		return err
	}

	keyFileData, err := os.ReadFile(keyFile)
	if err != nil {
		log.Fatalf("Key file read error: %s", err)
		return err
	}

	b32 := sha256.Sum256(keyFileData)
	key := b32[:]
	res.Block, err = aes.NewCipher(key)
	if err != nil {
		log.Fatalf("Error on aes.NewCipher: %s", err)
		return err
	}

	inputFile, err := filepath.Abs(res.source)
	if err != nil {
		log.Fatalf("Absolute path for file not found: %s", err)
		return err
	}

	res.Data, err = os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Source file read error: %s", err)
		return err
	}

	return nil
}

func (res *shResource) FileSafe() error {
	fname := filepath.Base(res.GetSource() + ".shc")
	if res.GetOut() != "" {
		var err error
		fname, err = filepath.Abs(res.GetOut())
		if err != nil {
			log.Fatalf("Absolute path for out file not found: %s", err)
			return err
		}
	} else {
		if res.Type == "decrypt" {
			fname = filepath.Base(strings.Replace(string(res.GetSource()), ".shc", "", 1))
		}
	}

	err := os.WriteFile(fname, res.Data, 0664)
	if err != nil {
		fmt.Println("Error write data to file", fname, err)
		os.Exit(0)
	}

	return nil
}

func (res *shResource) Encrypt() error {
	res.cbcBlockPadding()
	res.Data = append(make([]byte, aes.BlockSize), res.Data...)
	iv := res.Data[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatalf("Error while generating IV: %s", err)
		return err
	}

	mode := cipher.NewCBCEncrypter(res.Block, iv)
	mode.CryptBlocks(res.Data[aes.BlockSize:], res.Data[aes.BlockSize:])
	res.Data = append(res.Data, byte(res.padding))

	return nil
}

func (res *shResource) Decrypt() error {
	iv := res.Data[:aes.BlockSize]
	res.Data = res.Data[aes.BlockSize:]
	res.padding = int(res.Data[len(res.Data)-1])
	res.Data = res.Data[:len(res.Data)-1]
	mode := cipher.NewCBCDecrypter(res.Block, iv)
	mode.CryptBlocks(res.Data, res.Data)
	res.Data = res.Data[:len(res.Data)-res.padding]

	return nil
}

func (res *shResource) cbcBlockPadding() {
	res.padding = aes.BlockSize - len(res.Data)%aes.BlockSize

	if res.padding > 0 {
		paddBuf := bytes.Repeat([]byte{byte(0)}, res.padding)
		res.Data = append(res.Data, paddBuf...)
	}
}
