package models

import (
	"encoding/xml"
	"gopkg.in/validator.v2"
	"utils"
)

type Category struct {
	XMLName xml.Name `xml:"Category" json:"-"`
	Key     string   `datastore:"-"`
	Name    string   `validate:"nonzero"`
}

func (c *Category) Validate() error {
	return validator.Validate(c)
}

func (c *Category) Hash() string {
	return utils.Hash(c.Name)
}
