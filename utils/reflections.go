package utils

import "reflect"

func walkStructWithTag(v reflect.Value, tagPath []reflect.StructTag, f func(value reflect.Value, tagPath []reflect.StructTag) error) error {
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			tag := t.Field(i).Tag
			tagPath = append(tagPath, tag)
			err := walkStructWithTag(field, tagPath, f)
			if err != nil {
				return err
			}
			tagPath = tagPath[:len(tagPath)-1]
		}
	case reflect.Interface:
		return walkStructWithTag(v.Elem(), tagPath, f)
	default:
		return f(v, tagPath)
	}
	return nil
}

func WalkStructWithTag(v any, f func(value reflect.Value, tagPath []reflect.StructTag) error) error {
	if v == nil {
		return nil
	}
	return walkStructWithTag(reflect.Indirect(reflect.ValueOf(v)), make([]reflect.StructTag, 1), f)
}
