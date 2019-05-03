package util

import "reflect"

func ObjectIn(collect interface{}, obj interface{}) bool {
	slice := toSlice(collect)
	if collect == nil || len(slice) == 0 {
		return false
	}

	for _, v := range slice {
		if v == obj {
			return true
		}
	}

	return false
}

func toSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		panic("param not slice")
	}

	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}

	return ret
}
