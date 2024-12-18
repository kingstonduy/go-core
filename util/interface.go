package util

import "reflect"

func CheckNilInterface(i interface{}) bool {
	iv := reflect.ValueOf(i)
	if !iv.IsValid() {
		return true
	}
	switch iv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		return iv.IsNil()
	default:
		return false
	}
}
