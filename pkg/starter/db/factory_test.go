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

	db, err = factory.New("xxx")
	assert.Equal(t, "database is not implemented", err.Error())
}
