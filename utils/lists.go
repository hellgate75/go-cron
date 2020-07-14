package utils

import (
	"reflect"
	"strings"
)

func ListContains(target interface{}, list interface{}) (bool, int) {
	if reflect.TypeOf(list).Kind() == reflect.Slice || reflect.TypeOf(list).Kind() == reflect.Array {
		listValue := reflect.ValueOf(list)
		for i := 0; i < listValue.Len(); i++ {
			if target == listValue.Index(i).Interface() {
				return true, i
			}
		}
	}
	if reflect.TypeOf(target).Kind() == reflect.String && reflect.TypeOf(list).Kind() == reflect.String {
		return strings.Contains(list.(string), target.(string)), strings.Index(list.(string), target.(string))
	}
	return false, -1
}

func ListMatchContains(list []interface{}, match func(i interface{}) bool) (bool, int) {
	for idx, val := range list {
		if match(val) {
			return true, idx
		}
	}
	return false, -1
}