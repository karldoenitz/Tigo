package test_case

import (
	"Tigo/TigoWeb"
	"strings"
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

func TestMethodEnum(t *testing.T) {
	methods := []string{
		"get", "head", "put", "post", "delete", "connect", "options", "trace",
	}
	for _, method := range methods {
		m := strings.ToUpper(string(method[0])) + method[1:]
		if TigoWeb.MethodEnum(method) != m {
			t.Error("MethodEnum test failed")
			return
		}
		if TigoWeb.MethodEnum(strings.ToLower(method)) != m {
			t.Error("MethodEnum test failed")
			return
		}
	}
	t.Log("success")
}

func TestUrlEncode(t *testing.T) {
	originData := "测试用例1"
	encodedData := TigoWeb.UrlEncode(originData)
	decodedData := TigoWeb.UrlDecode(encodedData)
	if encodedData != decodedData {
		t.Error("url encode decode failed")
		return
	}
	t.Log("url encode decode testcase passed")
}
