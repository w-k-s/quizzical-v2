package api

import (
	"encoding/xml"
	"time"
)

type Response struct {
	Data interface{} `json:"Data"`
}

type Error struct {
	XMLName xml.Name `xml:"Error" json:"-"`
	Message string   `json:"Error" xml:"Message"`
}

func NewError(err error) *Error {
	return &Error{Message: err.Error()}
}

func NewErrorFromString(message string) *Error{
	return &Error{Message: message}
}

type Token struct{
	XMLName xml.Name `xml:"Token" json:"-"`
	Token string 
	ExpiryUnix int64
	ExpiryTimestamp string
}

func NewToken(token string, expiry time.Time) *Token{
	return &Token{
		Token: token,
		ExpiryUnix: expiry.Unix(),
		ExpiryTimestamp: expiry.Format(time.UnixDate),
	}
}