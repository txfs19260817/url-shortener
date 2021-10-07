package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitLogger(t *testing.T) {
	w := InitLogger()
	field := Log.WithField("test", "test")
	assert.Equal(t, "test", field.Data["test"])
	assert.IsType(t, Log.Formatter, &logrus.JSONFormatter{})
	assert.Equal(t, Log.GetLevel(), logrus.InfoLevel)
	assert.NoError(t, w.Close())
}
