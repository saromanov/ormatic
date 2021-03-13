package ormatic

import (
	"errors"
	"reflect"
	"strings"

	"github.com/saromanov/ormatic/models"
)

const dbField = "db"

var (
	errNoStruct = errors.New("provided data is not struct")
)

// getFieldsFromStruct returns list of fields with db tag
func getFieldsFromStruct(d interface{})([]models.Pair, error) {
	values := []models.Pair{}
	if ok := isStruct(d); !ok {
		return values, errNoStruct
	}
	val := reflect.ValueOf(d).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		dbTag := tag.Get(dbField)
		if dbTag == "" {
			dbTag = strings.ToLower(typeField.Name)
		}
		values = append(values, models.Pair{Key: dbTag, 
			Value:valueField.Interface(),
		})
	}
	return values, nil
}

func getObjectName(d interface{}) string {
	if ok := isStruct(d); !ok {
		return ""
	}
	t := reflect.TypeOf(d)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return strings.ToLower(t.Name())
}

func isStruct(d interface{}) bool {
	switch reflect.ValueOf(d).Kind() {
	case reflect.Struct:
		return true
	case reflect.Ptr:
		return reflect.ValueOf(d).Type().Elem().Kind() == reflect.Struct
	}
	return false
}

// Return 
func getStructFieldsTypes(d interface{}) error {
	if ok := isStruct(d); !ok {
		return errNoStruct
	}
	v := reflect.ValueOf(d).Elem()
    for j := 0; j < v.NumField(); j++ {
        f := v.Field(j)
        n := v.Type().Field(j).Name
        t := f.Type().String()
        fmt.Printf("Name: %s  Kind: %s  Type: %s\n", n, f.Kind(), t)
    }
	return nil
}