package WebFramework

import (
	"crypto/cipher"
	"crypto/des"
	"bytes"
)

// Des加密函数
func DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	cryptData := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化cryptData也可以
	blockMode.CryptBlocks(cryptData, origData)
	return cryptData, nil
}

// PKCS5填充函数
func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText) % blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// Des解密函数
func DesDecrypt(cryptData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(cryptData))
	blockMode.CryptBlocks(origData, cryptData)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

// PKCS5还原函数
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unPadding 次
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
