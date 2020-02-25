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

func GetFieldsFromNameAndSet(i interface{}, name string, v interface{}) {
	reflect.ValueOf(i).Elem().FieldByName(name).Set(reflect.ValueOf(v))
}
