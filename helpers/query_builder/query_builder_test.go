package query_builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildInsertQuery(t *testing.T) {
	rx := `INSERT INTO users \((.*)\) VALUES \(\$1, \$2, \$3\)`
	input := struct {
		Id      string   `db:"id"`
		Val     []string `db:"val"`
		Name    string   `db:"name"`
		Ignore1 string   `db:"-"`
		Ignore2 string
	}{
		Id:      "12345",
		Name:    "John",
		Val:     []string{"Hi"},
		Ignore1: "ignore1",
	}

	t.Run("Using a pointer to struct", func(t *testing.T) {
		query, values := BuildInsertQueryFromStruct("users", &input)
		assert.Regexp(t, rx, query)
		assert.NotNil(t, values)
	})
}
