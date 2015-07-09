package datastore

import (
	"appengine"
	"appengine/datastore"
	"models"
	"net/http"
)

const (
	EntityCategory       = "category"
	DefaultCategoryLimit = 30
)

type CategoryStore struct{}

func (store *CategoryStore) GetAll(request *http.Request, limit int) ([]*models.Category, error) {

	if limit <= 0 {
		limit = DefaultCategoryLimit
	}

	context := appengine.NewContext(request)

	query := datastore.NewQuery(EntityCategory).Limit(limit)

	var categories []*models.Category

	keys, err := query.GetAll(context, &categories)

	if err != nil {
		return nil, err
	}

	if categories == nil {
		categories = make([]*models.Category, 0)
	}

	for i, category := range categories {
		encodedKey := keys[i].Encode()
		category.Key = encodedKey
	}

	return categories, nil
}

func (s *CategoryStore) Save(request *http.Request, category *models.Category) error {

	context := appengine.NewContext(request)

	completeKey := datastore.NewKey(context, EntityCategory, category.Hash(), 0, nil)
	key, err := datastore.Put(context, completeKey, category)

	if err != nil {
		return err
	}

	category.Key = key.Encode()

	return nil
}

func (s *CategoryStore) Delete(request *http.Request, key string) error {

	context := appengine.NewContext(request)
	decodedKey, err := datastore.DecodeKey(key)

	if err != nil {
		return err
	}

	category := new(models.Category)

	err = datastore.Get(context, decodedKey, category)
	if err != nil {
		return err
	}

	err = datastore.Delete(context, decodedKey)
	if err != nil {
		return err
	}

	return nil
}
