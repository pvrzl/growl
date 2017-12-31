package growl

import (
	"os"
	"reflect"
	"unicode"
)

func IsFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func ToSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			if out[len(out)-1] != '_' {
				out = append(out, '_')
			}
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

func GetStructName(data interface{}) string {
	if t := reflect.TypeOf(data); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func GetValue(v reflect.Value) interface{} {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		return v.Elem().Interface()
	}

	return v.Interface()
}

func OpenConnectionStats() int {
	return connDb.DB().Stats().OpenConnections
}
