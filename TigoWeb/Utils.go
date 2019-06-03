package TigoWeb

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
)

//////////////////////////////////////////////////数据加密工具////////////////////////////////////////////////////////////

// Encrypt 方法用来根据key对原始数据进行加密，并将加密结果进行base64编码，
// 加密失败则返回空
//   - 此处以后会进行异常处理方面的优化
func Encrypt(src []byte, key []byte) string {
	encryptValue, _ := encrypt(src, key)
	return base64.StdEncoding.EncodeToString(encryptValue)
}

// Decrypt 方法会先对原始数据进行base64解码，然后根据key进行解密，
// 解密失败则返回空
//   - 此处以后会进行异常处理方面的优化
func Decrypt(src []byte, key []byte) []byte {
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

//////////////////////////////////////////////////初始化全局配置//////////////////////////////////////////////////////////

// InitGlobalConfig 方法用来初始化全局变量
func InitGlobalConfig(configPath string) {
	config := GlobalConfig{}
	config.Init(configPath)
	globalConfig = &config
}

// InitGlobalConfigWithObj 可使用TigoWeb.GlobalConfig的实例进行初始化全局变量
func InitGlobalConfigWithObj(config GlobalConfig)  {
	globalConfig = &config
}

//////////////////////////////////////////////////HTTP相关工具///////////////////////////////////////////////////////////

// getFormDataStr 获取报文体中的Form信息
// 将url.Values中的数据迭代取出，存入一个数组中，
// 并将字符串拼接成一个字符串
func getFormDataStr(form url.Values) string {
	var params []string
	for k, v := range form {
		value := ""
		if len(v) > 0 {
			value = v[0]
		}
		params = append(params, fmt.Sprintf("%s=%s", k, value))
	}
	return strings.Join(params, "&")
}
