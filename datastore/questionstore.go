package datastore

import (
	"appengine"
	"appengine/datastore"
	"models"
	"net/http"
)

const (
	EntityQuestion       = "question"
	QueryNoOffset        = -1
	QueryNoCategory      = ""
	DefaultQuestionLimit = 20
)

type QuestionStore struct{}

func (store *QuestionStore) GetAll(context appengine.Context, limit, offset int, category string) ([]*models.Question, error) {

	if limit <= 0 {
		limit = DefaultQuestionLimit
	}

	query := datastore.NewQuery(EntityQuestion)

	if len(category) > 0 {
		query = query.Filter("Category =", category)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	query = query.Limit(limit)

	var questions []*models.Question

	keys, err := query.GetAll(context, &questions)

	if err != nil {
		return nil, err
	}

	if questions == nil {
		questions = make([]*models.Question, 0)
	}

	for i, question := range questions {
		encodedKey := keys[i].Encode()
		question.Key = encodedKey
	}

	return questions, nil

}

func (s *QuestionStore) Count(request *http.Request, category string) (int, error) {

	context := appengine.NewContext(request)
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
