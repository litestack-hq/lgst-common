package query_builder

import (
	"fmt"
	"math"
	"reflect"
	"strings"
)

const (
	TAG_NAME = "db"
)

func buildInsertQuery(table string, input map[string]any, skipConflicting bool) (string, []any) {
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

	query += ")"

	if skipConflicting {
		query += " ON CONFLICT DO NOTHING"
	}

	query += " RETURNING *"

	return query, values
}

func BuildNextCursor(results []any, length int, fields []string) string {
	if results == nil {
		return ""
	}

	if len(results) <= length {
		return ""
	}

	nextItem := results[length]
	values := make([]string, len(fields))

	for i, field := range fields {
		values[i] = getTagValue(nextItem, field)
	}

	return strings.Join(values, ",")
}

func getTagValue(nextItem any, tagName string) string {
	nextItemValue := reflect.ValueOf(nextItem)
	if nextItemValue.Kind() == reflect.Ptr {
		nextItemValue = nextItemValue.Elem()
	}

	for i := 0; i < nextItemValue.NumField(); i++ {
		fieldName := nextItemValue.Type().Field(i).Tag.Get(TAG_NAME)
		field := nextItemValue.Field(i)
		kind := field.Kind()

		if fieldName != tagName || kind != reflect.String {
			continue
		}

		return nextItemValue.Field(i).String()
	}

	return ""
}

func BuildPaginationQueryFromModel(input PaginationQueryInput, model any) (string, []string) {
	cursorFields := []string{"uuid"}
	query := fmt.Sprintf("SELECT * FROM %s", input.Table)
	queryLimit := int(math.Abs(float64(input.Limit)) + 1)
	orderBy := ""

	if ok, tag := tagExists(TAG_NAME, input.Sort.Field, model); ok && input.Sort.Order.IsValid() {
		orderBy = fmt.Sprintf(
			" ORDER BY %s %s",
			input.Sort.Field,
			input.Sort.Order,
		)

		cursorFields = append(cursorFields, tag)
	} else {
		cursorFields = append(cursorFields, "created_at")
	}

	if input.NextCursor != "" {
		query = fmt.Sprintf(
			"SELECT * FROM %s WHERE (%s) > (%s)",
			input.Table,
			strings.Join(cursorFields, ","),
			input.NextCursor,
		)
	}

	query += fmt.Sprintf("%s LIMIT %d", orderBy, queryLimit)

	return query, cursorFields
}

func tagExists(tag string, value string, model any) (bool, string) {
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	for i := 0; i < modelValue.NumField(); i++ {
		fieldName := modelValue.Type().Field(i).Tag.Get(TAG_NAME)

		if fieldName == "" || fieldName == "-" {
			continue
		}

		if strings.ToLower(fieldName) == value {
			return true, fieldName
		}
	}

	return false, tag
}

func BuildInsertQueryFromModel(table string, model any, skipConflicting bool) (string, []any) {
	inputValues := map[string]any{}
	modelValue := reflect.ValueOf(model)

	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	for i := 0; i < modelValue.NumField(); i++ {
		tableFieldName := modelValue.Type().Field(i).Tag.Get(TAG_NAME)

		if tableFieldName == "" || tableFieldName == "-" {
			continue
		}

		field := modelValue.Field(i)
		kind := field.Kind()

		if kind == reflect.String || kind == reflect.Slice || kind == reflect.Map {
			if field.Len() == 0 {
				continue
			}
		}

		inputValues[tableFieldName] = modelValue.Field(i).Interface()
	}

	return buildInsertQuery(table, inputValues, skipConflicting)
}
