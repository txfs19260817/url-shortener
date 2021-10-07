package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMongoDB(t *testing.T) {
	var c MongoDBConfig
	db, err := NewMongoDB(c)
	assert.Nil(t, db)
	assert.Error(t, err)
}
