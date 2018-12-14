package binding

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

var (
	ErrStructure = errors.New("Unmarshal() expects struct input. ")
)

// Unmarshal 将url.Values转为struct
func Unmarshal(values url.Values, s interface{}) error {
	val := reflect.ValueOf(s)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ErrStructure
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return ErrStructure
	}
	return reflectValueFromTag(values, val)
}

func reflectValueFromTag(values url.Values, val reflect.Value) error {
	fieldType := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := fieldType.Field(i)
		tag := field.Tag.Get("form")
		if tag == "-" {
			continue
		}
		vField := val.Field(i)
		uv := getVal(values, tag)
		switch vField.Kind() {
		case reflect.String:
			vField.SetString(uv)
		case reflect.Bool:
			uv := parseValue(uv, "bool")
			b, err := strconv.ParseBool(uv)
			if err != nil {
				return fmt.Errorf("cast bool has error, expect type: %v ,val: %v ,query key: %v", vField.Type(), uv, tag)
			}
			vField.SetBool(b)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			uv := parseValue(uv, "uint")
			n, err := strconv.ParseUint(uv, 10, 64)
			if err != nil || vField.OverflowUint(n) {
				return fmt.Errorf("cast uint has error, expect type: %v ,val: %v ,query key: %v", vField.Type(), uv, tag)
			}
			vField.SetUint(n)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			uv := parseValue(uv, "uint")
			n, err := strconv.ParseInt(uv, 10, 64)
			if err != nil || vField.OverflowInt(n) {
				return fmt.Errorf("cast int has error, expect type: %v ,val: %v ,query key: %v", vField.Type(), uv, tag)
			}
			vField.SetInt(n)
		case reflect.Float32, reflect.Float64:
			uv := parseValue(uv, "uint")
			n, err := strconv.ParseFloat(uv, vField.Type().Bits())
			if err != nil || vField.OverflowFloat(n) {
				return fmt.Errorf("cast float has error, expect type: %v ,val: %v ,query key: %v", vField.Type(), uv, tag)
			}
			vField.SetFloat(n)
		default:
			return fmt.Errorf("unsupported type: %v ,val: %v ,query key: %v", vField.Type(), uv, tag)
		}
	}
	return nil
}

// parseValue 解析
func parseValue(value, tp string) string {
	if tp == "bool" && value == "" {
		return "false"
	} else if (tp == "int" || tp == "uint" || tp == "float") && value == "" {
		return "0"
	} else {
		return value
	}
}

// getVal 根据tag从url.Values中获取值
func getVal(values url.Values, tag string) string {
	return values.Get(tag)
}

// bytesToQuery 将bytes转换为url的values值
func bytesToQuery(urlParam []byte) (url.Values, error) {
	formatUrl := "http://www.query.com/param?"+string(urlParam)
	u, err := url.Parse(formatUrl)
	if err != nil {
		return nil, err
	}
	return url.ParseQuery(u.RawQuery)
}

// FormBytesToStructure 将x-www-form-urlencoded转换为struct实例
func FormBytesToStructure(form []byte, obj interface{}) error {
	values, err := bytesToQuery(form)
	if err != nil {
		return err
	}
	return Unmarshal(values, obj)
}

// ParseFormToInstance 将form转为structure对应的instance，并根据tag校验字段
func ParseFormToInstance(form []byte, obj interface{}) error {
	err := FormBytesToStructure(form, obj)
	if err != nil {
		return err
	}
	return ValidateInstance(obj)
}
