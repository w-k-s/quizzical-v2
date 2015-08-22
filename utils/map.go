package utils

import (
	"reflect"
	"strconv"
)

type MapHelper struct {
	Map map[string]interface{}
}

func (mapHelper MapHelper) Int(key string, defaultValue int) int {

	r := mapHelper.Map[key]

	if reflect.TypeOf(r).Kind() == reflect.Int ||
		reflect.TypeOf(r).Kind() == reflect.Uint ||
		reflect.TypeOf(r).Kind() == reflect.Float64 {
		return r.(int)
	}
	if reflect.TypeOf(r).Kind() == reflect.String{
		num,err := strconv.Atoi(r.(string))

		if err == nil {
			return num
		}

	}

	return defaultValue
}

func (mapHelper MapHelper) String(key string, defaultValue string) string {

	r := mapHelper.Map[key]

	if reflect.ValueOf(r).Kind() == reflect.String {
		return r.(string)
	}

	return defaultValue
}
