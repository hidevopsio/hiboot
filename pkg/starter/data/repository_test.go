package data

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBaseRepository(t *testing.T) {
	r := &BaseRepository{}

	err := r.DataSource()
	assert.Equal(t, NotImplemenedError, err)

	err = r.CloseDataSource()
	assert.Equal(t, NotImplemenedError, err)
}
