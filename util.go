package rms

import "reflect"

func containsArray(a interface{}, v interface{}) bool {
	if reflect.TypeOf(a).Kind() != reflect.Array {
		return false
	}
	if reflect.TypeOf(a).Elem().Kind() != reflect.TypeOf(v).Kind() {
		return false
	}
	i := 0
	arr := reflect.ValueOf(a)
	val := reflect.ValueOf(v)
	for i < arr.Len() {
		if val == arr.Index(i) {
			return true
		}
		i++
	}
	return false
}
