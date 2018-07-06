package TigoWeb

import (
	"encoding/base64"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"crypto/rand"
	"errors"
	"crypto/md5"
)

//////////////////////////////////////////////////数据加密工具////////////////////////////////////////////////////////////

// 根据key对原始数据进行加密，并将加密结果进行base64编码，
// 加密失败则返回空
//   - 此处以后会进行异常处理方面的优化
func Encrypt(src[]byte, key []byte) string {
	encryptValue, _ := encrypt(src, key)
	return base64.StdEncoding.EncodeToString(encryptValue)
}

// 先对原始数据进行base64解码，然后根据key进行解密，
// 解密失败则返回空
//   - 此处以后会进行异常处理方面的优化
func Decrypt(src[]byte, key []byte) ([]byte) {
	result, _ := base64.StdEncoding.DecodeString(string(src))
	value, _ := decrypt(result, key)
	return value
}

// aes加密函数，
// 先将key通过md5加密为64位，然后对原始值进行aes加密
func encrypt(plainText []byte, key []byte) ([]byte, error) {
	has := md5.Sum(key)
	hasKey := []byte(has[:])
	c, err := aes.NewCipher(hasKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plainText, nil), nil
}

// aes解密函数，
// 先将key通过md5加密为64位，然后对加密值进行aes解密
func decrypt(cipherText []byte, key []byte) ([]byte, error) {
	has := md5.Sum(key)
	hasKey := []byte(has[:])
	c, err := aes.NewCipher(hasKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, errors.New("cipherText too short")
	}
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	return gcm.Open(nil, nonce, cipherText, nil)
}
