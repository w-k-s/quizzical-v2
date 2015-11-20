package utils

import(
	"net/http"
	"strconv"
)

func FormInt(r * http.Request, key string, defaultValue int) int{

	if value,err := strconv.Atoi(r.FormValue(key)); err != nil{
		return defaultValue
	}else{
		return value
	}

}