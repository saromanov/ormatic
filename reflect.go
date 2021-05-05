package ormatic

import (
	"fmt"
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
// and parse of inner structure
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
		if valueField.Kind() == reflect.Struct {
			fields, err := getStructFieldsTypes(val.Field(i).Addr().Interface())
			if err != nil {
				return nil, errors.Wrap(err, "unable to get struct fields")
			}
			primary, tableName, err := getPrimaryKeyField(fields)
			if err != nil {
				return nil, errors.Wrap(err, "unable to get primary key from struct")
			}
			statement, err := getTagsFromRelationships(val.Type().Field(i).Tag, tableName, primary.Name)
			if err != nil {
				return nil, errors.Wrap(err, "unable to get tags from relationships")
			}

			child := reflect.ValueOf(val.Field(i).Addr().Interface()).Elem()
			primValue, err := getPrimaryKeyValue(child, statement.Source)
			if err != nil {
				return nil, errors.Wrap(err, "unable to get primary key value")
			}
			values = append(values, models.Pair{
				Key:   statement.Target,
				Value: primValue,
				Join: models.Join{
					Source: statement.Source,
					Target: statement.Target,
				},
			})
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

func getTagsFromRelationships(tags reflect.StructTag, tableName, key string) (models.Join, error) {
	dbTag := strings.ToLower(tags.Get("orm"))
	result := models.Join{}
	if !strings.Contains(dbTag, "on") {
		return result, nil
	}
	splitter := strings.Split(dbTag, "=")
	if len(splitter) != 2 {
		return result, nil
	}
	result.Target = tags.Get("db")
	result.Source = key
	return result, nil
}

func getPrimaryKeyValue(child reflect.Value, source interface{}) (interface{}, error) {
	for j := 0; j < child.NumField(); j++ {
		if strings.ToLower(child.Type().Field(j).Name) == source {
			return child.Field(j).Interface(), nil
		}
	}
	return nil, fmt.Errorf("unable to get primary key value")
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
				Name: getColumnName(v.Type().Field(j), v.Type().Field(j).Tag),
				Type: goTypeToSqlType[f.Type().String()],
				Tags: parseTableTags(v.Type().Field(j).Tag),
			})
		case reflect.Struct:
			tags := v.Type().Field(j).Tag.Get("orm")
			for _, t := range strings.Split(strings.ToLower(tags), ",") {
				if strings.Contains(t, "on") {
					res := strings.Split(t, "=")
					table, column, err := getTableAndColumnFromRels(res[1])
					if err != nil {
						return nil, err
					}
					root.Relationships = []models.Relationship{
						models.Relationship{
							TableName: table,
							Column:    column,
							Parent:    root.TableName,
							Name:      column,
						},
					}
				}
			}
			root.TableFields = append(root.TableFields, models.TableField{
				Name: getColumnName(v.Type().Field(j), v.Type().Field(j).Tag),
				Type: goTypeToSqlType["int"],
				Tags: parseTableTags(v.Type().Field(j).Tag),
			})

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

func getTableAndColumnFromRels(data string) (string, string, error) {
	result := strings.Split(data, ".")
	if len(result) != 2 {
		return "", "", errors.New("unable to parse relationship tag")
	}
	return result[0], result[1], nil
}

// return column name. If db tag is empty
// then return defined name
func getColumnName(sf reflect.StructField, tag reflect.StructTag) string {
	dbTag := tag.Get("db")
	if dbTag == "" {
		return strings.ToLower(sf.Name)
	}
	return dbTag
}

// return primary key from slice of fields
func getPrimaryKeyField(data []models.Create) (models.TableField, string, error) {
	for _, d := range data {
		for _, f := range d.TableFields {
			if f.Tags.PrimaryKey {
				return f, d.TableName, nil
			}
		}
	}
	return models.TableField{}, "", errors.New("unable to find primary key")
}
