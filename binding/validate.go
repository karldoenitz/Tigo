package binding

import (
	"errors"
	"fmt"
	"github.com/karldoenitz/Tigo/logger"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// checkField 对字段类型进行校验
func checkField(field reflect.StructField, vField reflect.Value) error {
	fieldKind := vField.Kind()
	switch fieldKind {
	case reflect.Ptr:
		if err := checkField(field, vField.Elem()); err != nil {
			return err
		}
	case reflect.Interface:
		if !vField.IsNil() {
			if err := checkField(field, vField.Elem()); err != nil {
				return err
			}
		}
	case reflect.Struct:
		if err := checkStructureField(field, vField); err != nil {
			return err
		}
	case reflect.Slice:
		if err := checkSliceField(field, vField); err != nil {
			return err
		}
	case reflect.Array:
		if err := checkArrayField(field, vField); err != nil {
			return err
		}
	case reflect.Map:
		if err := checkMapField(field, vField); err != nil {
			return err
		}
	case reflect.Invalid:
		if err := checkInvalidField(field, vField); err != nil {
			return err
		}
	case reflect.Uintptr, reflect.UnsafePointer, reflect.Chan, reflect.Func:
		logger.Warning.Printf("%s's kind is: %s", field.Name, fieldKind)
		logger.Warning.Printf("%s is unsupported field kind", field.Name)
		break
	default:
		if err := checkBasicField(field, vField); err != nil {
			return err
		}
	}
	return nil
}

// checkInvalidField 对无效字段进行校验
func checkInvalidField(field reflect.StructField, vField reflect.Value) error {
	required, isRequiredExisted := field.Tag.Lookup("required")
	if !isRequiredExisted || strings.ToLower(required) != "true" {
		return nil
	}
	defaultValue, isDefaultExisted := field.Tag.Lookup("default")
	regexStr, isRegexExisted := field.Tag.Lookup("regex")
	if isDefaultExisted {
		logger.Error.Printf("default value `%s` is invalid for nil field", defaultValue)
	}
	if isRegexExisted {
		logger.Error.Printf("regex `%s` is invalid for nil field", regexStr)
	}
	return fmt.Errorf("%s is a required field, can not be nil, value: %s", field.Name, vField.String())
}

// checkMapField 对Map类型的字段进行校验
func checkMapField(field reflect.StructField, vField reflect.Value) error {
	logger.Warning.Printf("Do not support map kind field: %s value: %s", field.Name, vField.String())
	return nil
}

// checkSliceField 对切片类型的字段进行校验
func checkSliceField(field reflect.StructField, vField reflect.Value) error {
	for i := 0; i < vField.Len(); i++ {
		v := vField.Index(i)
		t := v.Type()
		if v.Kind() != reflect.Struct && v.Kind() != reflect.Interface {
			logger.Warning.Printf("Only support interface/struct kind field: %s", field.Name)
			break
		}
		if err := checkObjBinding(t, v); err != nil {
			return errors.New(fmt.Sprintf("Field %s=>index %d has an error: %s", field.Name, i, err.Error()))
		}
	}
	return nil
}

// checkArrayField 对数组类型的字段进行校验
func checkArrayField(field reflect.StructField, vField reflect.Value) error {
	logger.Warning.Printf("Do not support array kind field: %s value: %s", field.Name, vField.String())
	return nil
}

// checkStructureField 对结构体类型的字段进行校验
func checkStructureField(field reflect.StructField, vField reflect.Value) error {
	required, isRequiredExisted := field.Tag.Lookup("required")
	if !isRequiredExisted || strings.ToLower(required) != "true" {
		return nil
	}
	if !(vField.CanAddr() && vField.Addr().CanInterface()) {
		return fmt.Errorf("%s's address can not be obtained with Addr", field.Name)
	}
	attrType := vField.Type()
	attrValue := vField
	for i := 0; i < attrType.NumField(); i++ {
		if err := checkField(attrType.Field(i), attrValue.Field(i)); err != nil {
			return err
		}
	}
	return nil
}

// checkBasicField 对基本数据类型的字段进行校验
func checkBasicField(field reflect.StructField, vField reflect.Value) error {
	var isRequired bool
	required, isRequiredExisted := field.Tag.Lookup("required")
	if isRequiredExisted && strings.ToLower(required) == "true" {
		isRequired = true
	}
	defaultValue, isDefaultExisted := field.Tag.Lookup("default")
	regexStr, isRegexExisted := field.Tag.Lookup("regex")
	fieldKind := vField.Type().Kind()
	fieldValue := vField.Interface()
	fieldName := field.Name
	switch fieldKind {
	case reflect.Bool:
		break
	case reflect.String:
		value := fieldValue.(string)
		if isRequired && value == "" && (!isDefaultExisted || defaultValue == "") {
			return RequiredErr(fieldName)
		}
		if value == "" && defaultValue != "" {
			value = defaultValue
		}
		if isRegexExisted && !isMatchRegex(value, regexStr) {
			return RegexErr(fieldName)
		}
		vField.SetString(value)
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, e := strconv.ParseInt(fmt.Sprint(fieldValue), 10, 64)
		if e != nil {
			return e
		}
		if isRequired && value == 0 && (!isDefaultExisted || defaultValue == "") {
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
		vField.SetInt(value)
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, e := strconv.ParseUint(fmt.Sprint(fieldValue), 10, 64)
		if e != nil {
			return e
		}
		if isRequired && value == 0 && (!isDefaultExisted || defaultValue == "") {
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
		vField.SetUint(value)
		break
	case reflect.Float32, reflect.Float64:
		value, e := strconv.ParseFloat(fmt.Sprint(fieldValue), 64)
		if e != nil {
			return e
		}
		if isRequired && value == 0 && (!isDefaultExisted || defaultValue == "") {
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
		vField.SetFloat(value)
		break
	}
	return nil
}

// checkObjBinding 检查实例的字段
func checkObjBinding(element reflect.Type, vElement reflect.Value) error {
	for i := 0; i < element.NumField(); i++ {
		if err := checkField(element.Field(i), vElement.Field(i)); err != nil {
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
	if err := checkObjBinding(element, vElement); err != nil {
		return err
	}
	return bindingCheck(v)
}

// 判断是否符合正则规则
func isMatchRegex(str, regex string) bool {
	rgx := regexp.MustCompile(regex)
	return rgx.MatchString(str)
}

// bindingCheck 调用param结构体的Check函数对param进行校验
func bindingCheck(vElement reflect.Value) (err error) {
	check := vElement.MethodByName("Check")
	var functionParams []reflect.Value
	if check.IsValid() {
		values := check.Call(functionParams)
		if !values[0].IsNil() {
			return values[0].Interface().(error)
		}
	}
	return
}
