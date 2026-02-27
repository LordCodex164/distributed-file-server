package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
)

func generateID() string {
	buf := make([]byte, 32)
	io.ReadFull(rand.Reader, buf)
	return hex.EncodeToString(buf)
}

func CopyStream(stream cipher.Stream) {

}

func copyEncrypt(key []byte, src io.Reader) (int, error) {
	//create a ciphertext
	log.Println("block")
	block, err := aes.NewCipher(key)
	log.Println("block", err)
	if err != nil {
		return 0, err
	}

	log.Println("block", block)

	//iv := make([]byte, block.BlockSize()) // 16 bytes

	//read the random text into byte slice
	iv := make([]byte, block.BlockSize()) // 16 bytes
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err
	}

	fmt.Println("block", iv)

	return 1, nil
	
}