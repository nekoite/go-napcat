package utils

import "reflect"

func DerefAny(v any) any {
	return reflect.Indirect(reflect.ValueOf(v)).Interface()
}
