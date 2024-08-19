package utils

import "reflect"

func walkStructWithTag(v reflect.Value, tagPath []reflect.StructTag, f func(value reflect.Value, tagPath []reflect.StructTag)) {
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			tag := t.Field(i).Tag
			tagPath = append(tagPath, tag)
			walkStructWithTag(field, tagPath, f)
			tagPath = tagPath[:len(tagPath)-1]
		}
	case reflect.Interface:
		walkStructWithTag(v.Elem(), tagPath, f)
	default:
		f(v, tagPath)
	}
}

func WalkStructWithTag(v any, f func(value reflect.Value, tagPath []reflect.StructTag)) {
	if v == nil {
		return
	}
	walkStructWithTag(reflect.Indirect(reflect.ValueOf(v)), make([]reflect.StructTag, 1), f)
}
