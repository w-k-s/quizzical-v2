package utils

import (
	"github.com/martini-contrib/sessions"
)

func PushFlash(session sessions.Session, key string, value interface{}) bool {

	//sanity check
	if len(key) == 0 || value == nil || session == nil {
		return false
	}

	session.Set(key, value)
	return true
}

func PopFlash(session sessions.Session, key string) interface{} {

	//sanity check
	if len(key) == 0 || session == nil {
		return ""
	}

	value := session.Get(key)
	session.Delete(key)

	return value
}
