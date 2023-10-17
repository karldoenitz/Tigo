package test_case

import (
	"Tigo/TigoWeb"
	"reflect"
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
	if originData != decodedData {
		t.Error("url encode decode failed")
		return
	}
	t.Log("url encode decode testcase passed")
}

var testIns int

type TestcaseStruct struct {
	// 测试反向映射调用结构体的函数
	Element int
}

func (p *TestcaseStruct) Test1(param int) {
	testIns = param
}

func (p *TestcaseStruct) Test2(param int) {
	p.Element = param
}

func TestVoidFuncCall(t *testing.T) {
	param := 9
	ts := TestcaseStruct{}
	nts := reflect.New(reflect.TypeOf(ts))
	TigoWeb.VoidFuncCall(nts, "Test1", reflect.ValueOf(param))
	if param != testIns {
		t.Error("testcase1 VoidFuncCall failed")
		return
	}
	TigoWeb.VoidFuncCall(nts, "Test2", reflect.ValueOf(param))
	if param != int(nts.Elem().FieldByName("Element").Int()) {
		t.Error("testcase2 VoidFuncCall failed")
		return
	}
	t.Log("success")
}
