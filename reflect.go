package ormatic

import (
	"reflect"
	"strings"

	"github.com/saromanov/ormatic/models"
)


// getFieldsFromStruct returns list of fields with db tag
func getFieldsFromStruct(d interface{})[]Pair {
	val := reflect.ValueOf(d).Elem()
	values := []models.Pair{}
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		dbTag := tag.Get("db")
		if dbTag == "" {
			dbTag = strings.ToLower(typeField.Name)
		}
		values = append(values, models.Pair{Key: dbTag, 
			Value:valueField.Interface(),
		})
	}
	return values
}