package models

import (
	"gopkg.in/validator.v2"
	"utils"
)

type Category struct {
	Key  string `datastore:"-"`
	Name string `validate:"nonzero"`
}

func (c *Category) Validate() error {
	return validator.Validate(c)
}

func (c *Category) Hash() string {
	return utils.Hash(c.Name)
}
