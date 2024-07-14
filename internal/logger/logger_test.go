package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	assert.NoError(t, InitLogger("info"))
	assert.Error(t, InitLogger("asdf"))
}
