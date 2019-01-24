// Package binding 字段校验功能包
// 提供了一个非常简单的json以及form的校验功能，支持多种类型的校验。
package binding

import (
	"encoding/json"
)

// ParseJsonToInstance 将json转为structure对应的instance，并根据tag校验字段
func ParseJsonToInstance(jsonBytes []byte, obj interface{}) error {
	err := json.Unmarshal(jsonBytes, obj)
	if err != nil {
		return err
	}
	return ValidateInstance(obj)
}
