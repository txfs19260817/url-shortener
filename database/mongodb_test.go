package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMongoDB(t *testing.T) {
	db, err := NewMongoDB("", "", "", 0)
	assert.Nil(t, db)
	assert.Error(t, err)
}
