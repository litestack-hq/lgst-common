package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Add tests to verify that default envs are bound
func TestNew(t *testing.T) {
	conf := New()

	assert.NotEmpty(t, conf.APP_NAME)
	assert.Equal(t, conf.APP_ENV, "testing")
}
