package tools

import "reflect"

func GetStructTag(f reflect.StructField, tagName string) string {
	return string(f.Tag.Get(tagName))
}
