package models

import (
	"encoding/xml"
	"utils"
)

type Question struct {
	XMLName  xml.Name `xml:"question" json:"-" datastore:"-"`
	Key      string   `xml:"key" json:"key" datastore:"-"`
	Question string   `xml:"ask" form:"question" binding:"required"`
	Answer   string   `xml:"correct,attr" datastore:",noindex" form:"answer" binding:"required"`
	Category string   `xml:"-" form:"category" binding:"required"`
	A        string   `xml:"A" datastore:",noindex" form:"a" binding:"required"`
	B        string   `xml:"B" datastore:",noindex" form:"b" binding:"required"`
	C        string   `xml:"C" datastore:",noindex" form:"c" binding:"required"`
	D        string   `xml:"D" datastore:",noindex" form:"d" binding:"required"`
}

func (q *Question) Hash() string {
	return utils.Hash(q.Question, q.Answer, q.Category)
}

type Questions struct {
	XMLName   xml.Name `xml:"Questions" json:"-"`
	Category  string   `xml:"Title,attr"`
	Questions []*Question
}
