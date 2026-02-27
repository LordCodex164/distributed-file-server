package main

import (
	"bytes"
	"testing"
)

func TestCryptoEncrypt(t *testing.T) {
	byteKey := []byte("new messages")
	var fileBuffer = new(bytes.Buffer)
	n, _ := copyEncrypt(byteKey, fileBuffer)
	t.Log(">>n", n)
}