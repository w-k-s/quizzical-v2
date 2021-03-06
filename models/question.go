package models

import (
	"encoding/xml"
	"gopkg.in/validator.v2"
	"utils"
)

type Question struct {
	XMLName  xml.Name `xml:"Question" json:"-"` //I made the mistake of not specifying datastore:"-" when I first created this model.
	Key      string   `datastore:"-"`
	Question string   `validate:"nonzero"`
	Answer   string   `datastore:",noindex" validate:"nonzero,regexp=[ABCD],min=1,max=1"`
	Category string   `validate:"nonzero"`
	A        string   `validate:"nonzero"`
	B        string   `validate:"nonzero"`
	C        string   `validate:"nonzero"`
	D        string   `validate:"nonzero"`
}

func (q *Question) Hash() string {
	return utils.Hash(q.Question, q.Answer, q.Category)
}

func (q *Question) Validate() error {
	return validator.Validate(q)
}
