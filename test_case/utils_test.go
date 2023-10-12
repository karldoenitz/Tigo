package test_case

import (
	"Tigo/TigoWeb"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	originData := "test data"
	key := "key"
	encryptData, err := TigoWeb.Encrypt([]byte(originData), []byte(key))
	if err != nil {
		t.Error(err.Error())
		return
	}
	decryptData, err := TigoWeb.Decrypt([]byte(encryptData), []byte(key))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if string(decryptData) != originData {
		t.Error("data invalid")
		return
	}
	t.Log("success")
}
