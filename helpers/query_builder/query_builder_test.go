package query_builder

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuildPaginationQueryFromModel(t *testing.T) {
	type User struct {
		Id        string    `db:"id" json:"id" fake:"{uuid}"`
		Active    bool      `db:"active" json:"active" fake:"{bool}"`
		Name      string    `db:"name" json:"name" fake:"{name}"`
		CreatedAt time.Time `db:"created_at" json:"created_at" fake:"skip"`
	}

	query, args := BuildPaginationQueryFromModel(PaginationQueryInput{
		Table: "users",
		Limit: 5,
	}, User{})

	assert.Equal(t, "SELECT * FROM users ORDER BY created_at ASC, id ASC LIMIT 6", query)
	assert.Equal(t, 0, len(args))

	query, args = BuildPaginationQueryFromModel(PaginationQueryInput{
		Table:      "users",
		Limit:      5,
		NextCursor: "MjAyMy0xMC0yOFQxODo1NDo1My41MjQxNTJaLDFkMjEzNDY1LTRjYzktNGI4Yy1hM2JmLWQ5MTFiODhiMTk3Nw==",
		Sort: struct {
			Field string
			Order TableSortOrder
		}{
			Field: "",
			Order: "ASC",
		},
	}, &User{})

	assert.Equal(t, "SELECT * FROM users WHERE (created_at, id) > ($1, $2) ORDER BY created_at ASC, id ASC LIMIT 6", query)
	assert.Equal(t, 2, len(args))

	query, args = BuildPaginationQueryFromModel(PaginationQueryInput{
		Table:      "users",
		Limit:      5,
		NextCursor: "MjAyMy0xMC0yOFQxODo1NDo1My41MjQxNTJaLDFkMjEzNDY1LTRjYzktNGI4Yy1hM2JmLWQ5MTFiODhiMTk3Nw==",
		Sort: struct {
			Field string
			Order TableSortOrder
		}{
			Field: "name",
			Order: "ASC",
		},
	}, &User{})

	assert.Equal(t, "SELECT * FROM users WHERE (created_at, id) > ($1, $2) ORDER BY created_at ASC, id ASC, name ASC LIMIT 6", query)
	assert.Equal(t, 2, len(args))
}

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
		query, values := BuildInsertQueryFromModel("users", input, false)
		queryIgnoringConflict, _ := BuildInsertQueryFromModel("users", input, true)

		assert.Regexp(t, regularQueryRx, query)
		assert.Regexp(t, safeQueryRx, queryIgnoringConflict)
		assert.NotNil(t, values)
	})

	t.Run("Using a pointer to struct", func(t *testing.T) {
		query, values := BuildInsertQueryFromModel("users", &input, false)
		assert.Regexp(t, regularQueryRx, query)
		assert.NotNil(t, values)
	})
}
