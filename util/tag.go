package util

import "reflect"

func GetFieldFromTag(i interface{}, tag, tagv string) string {
	fields := reflect.TypeOf(i).Elem()
	for i := 0; i < fields.NumField(); i++ {
		if fields.Field(i).Tag.Get(tag) == tagv {
			return fields.Field(i).Name
		}
	}
	return ""
}

func GetFieldTag(i interface{}, field, tagk string) string {
	fields := reflect.TypeOf(i).Elem()
	for i := 0; i < fields.NumField(); i++ {
		if fields.Field(i).Name == field {
			return fields.Field(i).Tag.Get(tagk)
		}
	}
	return ""
}

func GetFieldsFromNameAndSetInt(i interface{}, name string, v int) {
	reflect.ValueOf(i).Elem().FieldByName(name).SetInt(int64(v))
}

func GetFieldsFromNameAndSetString(i interface{}, name string, v string) {
	reflect.ValueOf(i).Elem().FieldByName(name).SetString(v)
}

func InvokeMapMethod(i interface{}, name string, v string) {
	param := []reflect.Value{reflect.ValueOf(v)}
	reflect.ValueOf(i).MethodByName(name).Call(param)
}
