package data

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	repo := new(BaseRepository)

	repo.SetName("foo")
	assert.Equal(t, "foo", repo.Name())

	repo.SetDataSource(nil)
	assert.Equal(t, nil, repo.DataSource())
}
