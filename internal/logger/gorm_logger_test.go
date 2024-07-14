package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_NewGormLogger(t *testing.T) {
	l := NewGormLogger(zap.NewNop(), GormLoggerConfig{})
	assert.NotNil(t, l)
}
