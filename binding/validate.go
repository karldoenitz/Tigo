package binding

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

// checkField 检查字段是否有效
func checkField(element reflect.Type, vElement reflect.Value, i int) error {
	var isRequired bool
	required, isRequiredExisted := element.Field(i).Tag.Lookup("required")
	if isRequiredExisted && (required == "true" || required == "TRUE") {
		isRequired = true
	}
	fieldType := element.Field(i).Type.Kind()
	fieldValue := vElement.Field(i).Interface()
	fieldName := element.Field(i).Name
	defaultValue, isDefaultExisted := element.Field(i).Tag.Lookup("default")
	regexStr, isRegexExisted := element.Field(i).Tag.Lookup("regex")
	switch fieldType {
	case reflect.Bool:
		break
	case reflect.String:
		value := fieldValue.(string)
		if isRequired && value == "" && (!isDefaultExisted || defaultValue=="") {
			return RequiredErr(fieldName)
		}
		if value == "" && defaultValue != "" {
			value = defaultValue
		}
		if isRegexExisted && !isMatchRegex(value, regexStr) {
			return RegexErr(fieldName)
		}
		vElement.Field(i).SetString(value)
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, e := strconv.ParseInt(fmt.Sprint(fieldValue), 10, 64)
		if e != nil {
			return e
		}
		if isRequired && value == 0 && (!isDefaultExisted || defaultValue=="") {
			return RequiredErr(fieldName)
		}
		if value == 0 && defaultValue != "" {
			var err error
			value, err = strconv.ParseInt(defaultValue, 10, 64)
			if err != nil {
				return DefaultErr(fieldName)
			}
		}
		if isRegexExisted {
			valueStr := strconv.FormatInt(int64(value), 10)
			if !isMatchRegex(valueStr, regexStr) {
				return RegexErr(fieldName)
			}
		}
		vElement.Field(i).SetInt(value)
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, e := strconv.ParseUint(fmt.Sprint(fieldValue), 10, 64)
		if e != nil {
			return e
		}
		if isRequired && value == 0 && (!isDefaultExisted || defaultValue=="") {
			return RequiredErr(fieldName)
		}
		if value == 0 && defaultValue != "" {
			val, err := strconv.ParseUint(defaultValue, 10, 64)
			if err != nil {
				return DefaultErr(fieldName)
			}
			value = val
		}
		if isRegexExisted {
			valueStr := strconv.FormatUint(value, 10)
			if !isMatchRegex(valueStr, regexStr) {
				return RegexErr(fieldName)
			}
		}
		vElement.Field(i).SetUint(value)
		break
	case reflect.Float32, reflect.Float64:
		value, e := strconv.ParseFloat(fmt.Sprint(fieldValue), 64)
		if e != nil {
			return e
		}
		if isRequired && value == 0 && (!isDefaultExisted || defaultValue=="") {
			return RequiredErr(fieldName)
		}
		if value == 0 && defaultValue != "" {
			val, err := strconv.ParseFloat(defaultValue, 64)
			if err != nil {
				return DefaultErr(fieldName)
			}
			value = val
		}
		if isRegexExisted {
			valueStr := strconv.FormatFloat(value, 'E', -1, 32)
			if !isMatchRegex(valueStr, regexStr) {
				return RegexErr(fieldName)
			}
		}
		vElement.Field(i).SetFloat(value)
		break
	default:
		field := vElement.Field(i)
		fieldKind := field.Kind()
		if fieldKind == reflect.Struct {
			if field.CanAddr() && field.Addr().CanInterface() {
				attrType := field.Type()
				attrValue := field
				return checkObjBinding(attrType, attrValue)
			}
		}
		if fieldKind == reflect.Ptr {
			if field.CanAddr() && field.Addr().CanInterface() {
				attrType := field.Type().Elem()
				attrValue := field.Elem()
				return checkObjBinding(attrType, attrValue)
			}
		}
	}
	return nil
}

// checkObjBinding 检查实例的字段
func checkObjBinding(element reflect.Type, vElement reflect.Value) error {
	for i := 0; i < element.NumField(); i++ {
		if err := checkField(element, vElement, i); err != nil {
			return err
		}
	}
	return nil
}

// ValidateInstance 检查结构体实例化是否有效
func ValidateInstance(obj interface{}) error {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	element := t.Elem()
	vElement := v.Elem()
	return checkObjBinding(element, vElement)
}

// 判断是否符合正则规则
func isMatchRegex(str, regex string) bool {
	rgx := regexp.MustCompile(regex)
	return rgx.MatchString(str)
}
