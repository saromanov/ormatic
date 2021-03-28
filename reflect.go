package ormatic

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/saromanov/ormatic/models"
)

const dbField = "db"

var (
	errNoStruct = errors.New("provided data is not struct")
)

var goTypeToSqlType = map[string]string{
	"int":     "integer",
	"int16":   "integer",
	"int32":   "integer",
	"int64":   "bigint",
	"string":  "text",
	"float32": "double",
	"float64": "double",
}

// getFieldsFromStruct returns list of fields with db tag
func getFieldsFromStruct(d interface{}) ([]models.Pair, error) {
	values := []models.Pair{}
	if ok := isStruct(d); !ok {
		return values, errNoStruct
	}
	val := reflect.ValueOf(d).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		if valueField.IsZero() || valueField.IsZero() {
			continue
		}
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		dbTag := tag.Get(dbField)
		if dbTag == "" {
			dbTag = strings.ToLower(typeField.Name)
		}
		values = append(values, models.Pair{Key: dbTag,
			Value: valueField.Interface(),
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

// Return struct for create table from the model
func getStructFieldsTypes(d interface{}) ([]models.Create, error) {
	resp := []models.Create{}
	if ok := isStruct(d); !ok {
		return nil, errNoStruct
	}
	v := reflect.ValueOf(d).Elem()
	root := models.Create{}
	root.TableName = getObjectName(d)
	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)
		switch f.Kind() {
		case reflect.String, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32,
			reflect.Float64:
			root.TableFields = append(root.TableFields, models.TableField{
				Name: strings.ToLower(v.Type().Field(j).Name),
				Type: goTypeToSqlType[f.Type().String()],
				Tags: parseTableTags(v.Type().Field(j).Tag),
			})
		case reflect.Struct:
			inner, err := getStructFieldsTypes(v.Field(j).Addr().Interface())
			if err != nil {
				return nil, errors.Wrap(err, "unable to get struct field")
			}
			resp = append(resp, inner...)
		}
	}
	resp = append(resp, root)
	return resp, nil
}

func parseTableTags(s reflect.StructTag) models.Tags {
	res := models.Tags{}
	tags := s.Get("orm")
	if tags == "" {
		return res
	}
	tags = strings.ToLower(tags)
	if strings.Contains(tags, "primary_key") {
		res.PrimaryKey = true
	}
	if strings.Contains(tags, "not_null") {
		res.NotNULL = true
	}
	if strings.Contains(tags, "unique") {
		res.Unique = true
	}
	if strings.Contains(tags, "index") {
		res.Index = "index"
	}
	return res
}
