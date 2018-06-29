package db

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {
	factory := new(DataSourceFactory)

	db, err := factory.New(DataSourceTypeBolt)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, db)
}
