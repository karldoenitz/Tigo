package binding

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
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
			b, err := strconv.ParseBool(uv)
			if err != nil {
				return fmt.Errorf("cast bool has error, expect type: %v ,val: %v ,query key: %v", vField.Type(), uv, tag)
			}
			vField.SetBool(b)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n, err := strconv.ParseUint(uv, 10, 64)
			if err != nil || vField.OverflowUint(n) {
				return fmt.Errorf("cast uint has error, expect type: %v ,val: %v ,query key: %v", vField.Type(), uv, tag)
			}
			vField.SetUint(n)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(uv, 10, 64)
			if err != nil || vField.OverflowInt(n) {
				return fmt.Errorf("cast int has error, expect type: %v ,val: %v ,query key: %v", vField.Type(), uv, tag)
			}
			vField.SetInt(n)
		case reflect.Float32, reflect.Float64:
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

// getVal 从tag中获取默认值，并在url.Values中没有值的时候设置默认值
func getVal(values url.Values, tag string) string {
	name, opts := parseTag(tag)
	uv := values.Get(name)
	optsLen := len(opts)
	if optsLen > 0 {
		if optsLen == 1 && uv == "" {
			uv = opts[0]
		}
	}
	return uv
}

type tagOptions []string

func parseTag(tag string) (string, tagOptions) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}
