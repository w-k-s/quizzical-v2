package models

import (
	"utils"
)

type Category struct {
	Key  string `datastore:"-"`
	Name string 
}

func (c *Category) Hash() string {
	return utils.Hash(c.Name)
}

type Categories struct {
	Categories []*Category `xml:"Category"`
}
