package query_builder

import (
	"fmt"
	"reflect"
)

func buildInsertQuery(table string, input map[string]any) (string, []any) {
	keys := reflect.ValueOf(input).MapKeys()
	values := []any{}
	fields := ""

	query := fmt.Sprintf("INSERT INTO %s (", table)
	for i, v := range keys {
		n := i + 1
		key := v.String()
		fields += key
		if n < len(keys) {
			fields += ", "
		}

		values = append(values, input[key])
	}

	query += fields + ") VALUES ("

	for i := range keys {
		n := i + 1
		query += fmt.Sprintf("$%d", n)
		if n < len(keys) {
			query += ", "
		}
	}

	query += ") RETURNING *"

	return query, values
}

func BuildInsertQueryFromStruct(table string, s any) (string, []any) {
	tagName := "db"
	input := map[string]any{}
	value := reflect.ValueOf(s)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

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
