package api

import (
	"encoding/xml"
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
