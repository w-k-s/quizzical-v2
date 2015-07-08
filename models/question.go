package models

import(
	"encoding/xml"
	"utils"
)

type Question struct {
	XMLName  xml.Name `xml:"question" json:"-" datastore:"-"`
	Key      string   `xml:"key" json:"key" datastore:"-"`
	Question string   `xml:"ask"`
	Answer   string   `xml:"correct,attr" datastore:",noindex"`
	Category string   `xml:"-"`
	A        string   `xml:"A" datastore:",noindex"`
	B        string   `xml:"B" datastore:",noindex"`
	C        string   `xml:"C" datastore:",noindex"`
	D        string   `xml:"D" datastore:",noindex"`
}

func (q * Question) Hash() string{
	return utils.Hash(q.Question,q.Answer,q.Category)
}

type Questions struct {
	XMLName   xml.Name `xml:"Questions" json:"-"`
	Category  string   `xml:"Title,attr"`
	Questions []*Question
}