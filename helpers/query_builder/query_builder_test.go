package query_builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildInsertQuery(t *testing.T) {
	rx := `INSERT INTO users \((.*)\) VALUES \(\$1, \$2, \$3, \$4\)`

	input := struct {
		Id            string   `db:"id"`
		Name          string   `db:"name"`
		StringSlice   []string `db:"val"`
		WalletBalance int      `db:"wallet_balance"`
		Ignore1       string   `db:"-"`
		Ignore2       string
	}{
		Id:            "12345",
		Name:          "John",
		StringSlice:   []string{"Hi"},
		WalletBalance: 500,
		Ignore1:       "ignore1",
	}

	output := []interface{}{
		"12345",
		"John",
		[]string{"Hi"},
		500,
	}

	t.Run("Using a struct value", func(t *testing.T) {
		query, values := BuildInsertQueryFromStruct("users", input)
		assert.Regexp(t, rx, query)
		assert.NotNil(t, values)
		assert.Equal(t, values, output)
	})

	t.Run("Using a pointer to struct", func(t *testing.T) {
		// NOTE: When a struct pointer is passed as input, the order
		// of the output slice might change after `value = value.Elem()`
		query, values := BuildInsertQueryFromStruct("users", &input)
		assert.Regexp(t, rx, query)
		assert.NotNil(t, values)
	})
}
