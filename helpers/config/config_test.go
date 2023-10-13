package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	conf := New()

	assert.NotEmpty(t, conf.APP_NAME)
	assert.Equal(t, conf.APP_ENV, "testing")
}
