package datastore

import (
	"appengine"
	"appengine/datastore"
	"reflect"
	"fmt"
)

type Repository struct{}

func (repo *Repository) GetAll(context appengine.Context, entityName string, items interface{},limit int) (error) {

	itemsValue := reflect.ValueOf(items).Elem()

	if itemsValue.Kind() != reflect.Slice{
		return fmt.Errorf("Expected Slice. Got %s",itemsValue.Kind())
	}

	if limit <= 0 {
		limit = DefaultCategoryLimit
	}

	query := datastore.NewQuery(entityName).Limit(limit)

	keys, err := query.GetAll(context, items)

	if err != nil {
		return err
	}

	for i := 0 ; i < itemsValue.Len() ; i++{
		itemValue := itemsValue.Index(i);
		if itemValue.Kind() == reflect.Struct {
			fieldValue := itemValue.FieldByName("Key")
			if fieldValue.CanSet() {
				fieldValue.SetString(keys[i].Encode())
			}
		}
	}

	return nil
}
