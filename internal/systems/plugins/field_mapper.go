package plugins

import (
	"reflect"
	"strings"
	"unicode"
)

type lowerCamelNameMapper struct{}

func (lowerCamelNameMapper) FieldName(_ reflect.Type, f reflect.StructField) string {
	return toLowerCamel(f.Name)
}

func (lowerCamelNameMapper) MethodName(_ reflect.Type, m reflect.Method) string {
	return toLowerCamel(m.Name)
}

func toLowerCamel(name string) string {
	if isUpper(name) {
		return strings.ToLower(name)
	} else {
		return strings.ToLower(name[:1]) + name[1:]
	}
}

func isUpper(s string) bool {
	for _, char := range s {
		if unicode.IsLower(char) {
			return false
		}
	}
	return true
}
