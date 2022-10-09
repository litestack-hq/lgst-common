package query_builder

import (
	"fmt"
	"reflect"
)

func buildInsertQuery(table string, input map[string]any) (string, []any) {
	keys := reflect.ValueOf(input).MapKeys()
	values := []any{}

	query := fmt.Sprintf("INSERT INTO %s (", table)
	for i, v := range keys {
		n := i + 1
		key := v.String()
		query += v.String()
		if n < len(keys) {
			query += ", "
		}

		values = append(values, input[key])
	}
	query += ") VALUES ("

	for i := range keys {
		n := i + 1
		query += fmt.Sprintf("$%d", n)
		if n < len(keys) {
			query += ", "
		}
	}
	query += ")"

	return query, values
}

func BuildInsertQueryFromStruct(table string, s any) (string, []any) {
	value := reflect.ValueOf(s)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	input := map[string]any{}
	tagName := "db"

	for i := 0; i < value.NumField(); i++ {
		tag := value.Type().Field(i).Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}

		field := value.Field(i)
		kind := field.Kind()

		if kind == reflect.String || kind == reflect.Slice || kind == reflect.Map {
			if field.Len() == 0 {
				continue
			}
		}

		input[tag] = value.Field(i).Interface()
	}

	return buildInsertQuery(table, input)
}
