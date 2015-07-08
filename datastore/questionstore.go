package datastore

import(
	"appengine"
	"appengine/datastore"
	"models"
	"net/http"
	"fmt"
	"math/rand"
)

const(
	EntityQuestion = "question"
	QueryNoOffset = -1
	QueryNoCategory = ""
)

type QuestionStore struct{}

func (store * QuestionStore) getQuestions(request *http.Request, limit,offset int, category string) ([] *models.Question,error){

	if limit <= 0{
		return nil,fmt.Errorf("Limit must be greater than 0.")
	}

	context := appengine.NewContext(request)

	query := datastore.NewQuery(EntityQuestion)

	if len(category) > 0 {
		query = query.Filter("Category =",category)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	query = query.Limit(limit)

	var questions []*models.Question

	keys,err := query.GetAll(context,&questions)

	if err != nil{
		return nil,err
	}

	if questions == nil {
		questions = make([]*models.Question,0)
	}

	for i,question := range questions {
		encodedKey := keys[i].Encode()
		question.Key = encodedKey
	}

	return questions,nil

}

func (store * QuestionStore) GetAll(request *http.Request, limit int) ([]*models.Question,error){

	return store.getQuestions(request,limit,QueryNoOffset,QueryNoCategory)	
}

func (store * QuestionStore) GetForCategory(request *http.Request,category string, limit int) ([]*models.Question,error){

	if len(category) == 0 {
		return nil,fmt.Errorf("category must not empty")
	}

	return store.getQuestions(request,limit,QueryNoOffset,category)
}

func (s * QuestionStore) Count(request * http.Request, category string) (int,error){

	context := appengine.NewContext(request)
	query := datastore.NewQuery(EntityQuestion)

	if len(category) > 0 {
		query = query.Filter("Category =",category)
	}

	count, err := query.Count(context)

	if err != nil {
		return -1, err
	}

	return count, nil
}

func (s * QuestionStore) Random(request * http.Request, category string, limit int) ([]*models.Question,error){

	if len(category) == 0 {
		return nil,fmt.Errorf("category must not empty")
	}

	count,err := s.Count(request,category)

	if err != nil{
		return nil,err
	}

	if limit >= count {
		
		return s.GetForCategory(request,category,limit)

	}else{

		/*
		If there are 30 questions, and we wish to deliver 10,
		then the offser should be such that offset + limit <= count.
		Therefore, maxOffset = count = limit
		*/
		maxOffset := count - limit
		offset := int(rand.Int63n(int64(maxOffset)))

		return s.getQuestions(request,limit,offset,category)
	}

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

func (s *QuestionStore) Save(request *http.Request, question *models.Question) error {
	
	context := appengine.NewContext(request)
	
	completeKey := datastore.NewKey(context, EntityQuestion,question.Hash(),0, nil)
	key, err := datastore.Put(context, completeKey, question)
	
	if err != nil {
		return err
	}
	
	question.Key = key.Encode()
	return nil
}

func (s *QuestionStore) Delete(request *http.Request, key string) error {
	
	context := appengine.NewContext(request)
	decodedKey, err := datastore.DecodeKey(key)
	
	if err != nil {
		return err
	}
	
	question := new(models.Question)
	
	err = datastore.Get(context, decodedKey, question)
	if err != nil {
		return err
	}

	err = datastore.Delete(context, decodedKey)
	if err != nil {
		return err
	}

	return nil
}