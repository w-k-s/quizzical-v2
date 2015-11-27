package api

import (
	"encoding/xml"
	"time"
	"math"
)

type Response struct {
	Data interface{} `json:"Data"`
}

type PaginatedResponse struct{
	Data interface{}
	CurrentPage int
	TotalPages int
	PageSize int
}

func NewPaginatedResponse(data interface{},pageSize, pageNumber, totalItems int) *PaginatedResponse{

	totalPages := 0
	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(totalItems)/float64(pageSize)))
	}

	return &PaginatedResponse{
		Data: data,
		CurrentPage: pageNumber,
		PageSize: pageSize,
		TotalPages: totalPages,
	}

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