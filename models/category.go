package models

import(
	"utils"
)

type Category struct {
    Key  string   `xml:"key" json:"key" datastore:"-"`
	Name string 
}

func (c * Category) Hash() string{
	return utils.Hash(c.Name)
}

type Categories struct{
	Categories []*Category `xml:"Category"`
}