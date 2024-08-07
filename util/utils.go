package util

import (
	"fmt"
	"reflect"
)

func GetVariableValueByName(variableName string, data interface{}) (interface{}, error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("provided data is not a struct")
	}

	field := v.FieldByName(variableName)
	if !field.IsValid() {
		return nil, fmt.Errorf("no such field: %s in data", variableName)
	}

	return field.Interface(), nil
}
