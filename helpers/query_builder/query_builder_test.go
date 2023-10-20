package query_builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildInsertQuery(t *testing.T) {
	regularQueryRx := `INSERT INTO users \((\b.*\b,\s){3}(\b.*\b)\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING \*`
	safeQueryRx := `INSERT INTO users \((\b.*\b,\s){3}(\b.*\b)\) VALUES \(\$1, \$2, \$3, \$4\) ON CONFLICT DO NOTHING RETURNING \*`

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

	// NOTE: The order of the output slice might change.
	// Asserting that will make the test flaky.
	t.Run("Using a struct value", func(t *testing.T) {
		query, values := BuildInsertQueryFromStruct("users", input, false)
		queryIgnoringConflict, _ := BuildInsertQueryFromStruct("users", input, true)

		assert.Regexp(t, regularQueryRx, query)
		assert.Regexp(t, safeQueryRx, queryIgnoringConflict)
		assert.NotNil(t, values)
	})

	t.Run("Using a pointer to struct", func(t *testing.T) {
		query, values := BuildInsertQueryFromStruct("users", &input, false)
		assert.Regexp(t, regularQueryRx, query)
		assert.NotNil(t, values)
	})
}
