package query_builder

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"math"
	"reflect"
	"strings"
	"time"
)

const (
	TAG_NAME = "db"
)

func BuildPaginationQueryFromModel(input PaginationQueryInput, model any) (string, []any) {
	queryLimit := 1
	query := fmt.Sprintf("SELECT * FROM %s", input.Table)
	args := []any{}
	queryLimit = int(queryLimit + int(math.Abs(float64(input.Limit))))
	orderBy := " ORDER BY created_at ASC, id ASC"

	useCustomSorting, sortFieldName := getFieldNameIfExists(TAG_NAME, input.Sort.Field, model)
	useCustomSorting = useCustomSorting && input.Sort.Order.IsValid()

	if useCustomSorting {
		orderBy = fmt.Sprintf(
			" ORDER BY %s %s, id ASC",
			sortFieldName,
			input.Sort.Order,
		)
	}

	if input.NextCursor != "" {
		decodedBytes, err := base64.StdEncoding.DecodeString(input.NextCursor)
		if err != nil {
			slog.Error(
				"failed to decode cursor",
				"value", input.NextCursor,
				"error", err,
			)
		}

		cursor := strings.Split(string(decodedBytes), ",")
		if len(cursor) == 2 {
			query = fmt.Sprintf("SELECT * FROM %s WHERE id > $1", input.Table)

			if !useCustomSorting {
				query = fmt.Sprintf("SELECT * FROM %s WHERE (created_at, id) > ($1, $2)", input.Table)
				parsedTime, err := time.Parse(time.RFC3339Nano, cursor[0])

				if err != nil {
					slog.Error(
						"failed to parse cursor created_at",
						"value", cursor[0],
						"format", time.RFC3339Nano,
						"error", err,
					)
				}

				args = append(args, parsedTime)
			}

			args = append(args, cursor[1])
		}
	}

	query += fmt.Sprintf("%s LIMIT %d", orderBy, queryLimit)

	return query, args
}

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

func getFieldNameIfExists(_ string, value string, model any) (bool, string) {
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

	return false, ""
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
