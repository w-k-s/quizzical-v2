package datastore

import (
	"appengine"
	"appengine/datastore"
	"models"
	"fmt"
)

const (
	EntityQuestion       = "question"
)

type QuestionStore struct{}

func (store *QuestionStore) GetAll(context appengine.Context, pageSize, pageNumber int, category string) (Result, error) {

	if len(category) == 0{
		return Result{}, fmt.Errorf("Category is required in order to query questions.")
	}

	if pageSize <= 0 {
		return Result{}, fmt.Errorf("Page Size must be > 0")
	}

	count,_ := store.Count(context,category);

	query := datastore.NewQuery(EntityQuestion).
		Filter("Category =", category)
	

	if pageNumber > 0 {

		offset := pageSize * (pageNumber-1)
		query = query.Offset(offset)
	}

	query = query.Limit(pageSize)

	var questions []*models.Question

	keys, err := query.GetAll(context, &questions)

	if err == nil {
		for i, question := range questions {
			encodedKey := keys[i].Encode()
			question.Key = encodedKey
		}
	}

	return Result{TotalCount: count, Data: questions}, err
}

func (s *QuestionStore) Count(context appengine.Context, category string) (int, error) {

	query := datastore.NewQuery(EntityQuestion)

	if len(category) > 0 {
		query = query.Filter("Category =", category)
	}

	return query.Count(context)
}

func (s *QuestionStore) Save(context appengine.Context, question *models.Question) error {

	completeKey := datastore.NewKey(context, EntityQuestion, question.Hash(), 0, nil)
	key, err := datastore.Put(context, completeKey, question)

	if err == nil {
		question.Key = key.Encode()
	}

	return err
}

func (s *QuestionStore) Delete(context appengine.Context, key string) error {

	decodedKey, err := datastore.DecodeKey(key)

	if err == nil {
		err = datastore.Delete(context, decodedKey)
	}

	return err
}
