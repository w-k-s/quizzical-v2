package utils

import(
	"net/http"
	"strconv"
)

func FormUInt(r * http.Request, key string, defaultValue int) int{

	if value,err := strconv.Atoi(r.FormValue(key)); err != nil || value < 0{
		return defaultValue
	}else{
		return value
	}

}