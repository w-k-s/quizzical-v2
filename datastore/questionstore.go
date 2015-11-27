package datastore

import (
	"appengine"
	"appengine/datastore"
	"models"
	"fmt"
	"net/http"
)

const (
	EntityQuestion       = "question"
	QueryNoOffset        = -1
	QueryNoCategory      = ""
	DefaultQuestionPageSize = 20
)

type QuestionStore struct{}

func (store *QuestionStore) GetAll(context appengine.Context, pageSize, pageNumber int, category string) (Result, error) {

	if len(category) == 0{
		return Result{}, fmt.Errorf("Category is required in order to query questions.")
	}

	if pageSize <= 0 {
		pageSize = DefaultQuestionPageSize
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

	if err != nil {
		return Result{}, err
	}

	for i, question := range questions {
		encodedKey := keys[i].Encode()
		question.Key = encodedKey
	}

	return Result{TotalCount: count, Data: questions}, nil
}

func (s *QuestionStore) Count(context appengine.Context, category string) (int, error) {

	query := datastore.NewQuery(EntityQuestion)

	if len(category) > 0 {
		query = query.Filter("Category =", category)
	}

	count, err := query.Count(context)

	if err != nil {
		return -1, err
	}

	return count, nil
}

func (s *QuestionStore) Find(request *http.Request, key string) (*models.Question, error) {

	context := appengine.NewContext(request)

	decodedKey, err := datastore.DecodeKey(key)
	if err != nil {
		return nil, err
	}

	question := new(models.Question)

	err = datastore.Get(context, decodedKey, question)
	if err != nil {
		return nil, err
	}

	question.Key = key
	return question, nil
}

func (s *QuestionStore) Save(context appengine.Context, question *models.Question) error {

	completeKey := datastore.NewKey(context, EntityQuestion, question.Hash(), 0, nil)
	key, err := datastore.Put(context, completeKey, question)

	if err != nil {
		return err
	}

	question.Key = key.Encode()

	return nil
}

func (s *QuestionStore) Delete(context appengine.Context, key string) error {

	decodedKey, err := datastore.DecodeKey(key)

	if err == nil {
		err = datastore.Delete(context, decodedKey)
	}

	return err
}
