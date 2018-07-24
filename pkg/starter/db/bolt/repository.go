package bolt

import (
	"errors"
	"github.com/hidevopsio/hiboot/pkg/starter/db"
)

type Repository interface {
	Put(key, value []byte) error
	Get(key []byte) (result []byte, err error)
	Delete(key []byte) (err error)
}

type repository struct {
	db.BaseRepository
}

var NilPointError = errors.New("dataSource is nil")

func (r *repository) ensureDataSource() (*bolt, error) {
	dataSource := r.DataSource()
	if dataSource == nil {
		return nil, NilPointError
	}
	return dataSource.(*bolt), nil
}

// Put inserts a key:value pair into the database
func (r *repository) Put(key, value []byte) error {
	bolt, err := r.ensureDataSource()
	if err != nil {
		return err
	}
	return bolt.Put([]byte(r.Name()), key, value)
}

// Get retrieves a key:value pair from the database
func (r *repository) Get(key []byte) ([]byte, error)  {
	bolt, err := r.ensureDataSource()
	if err != nil {
		return nil, err
	}
	return bolt.Get([]byte(r.Name()), key)
}

// Delete removes a key:value pair from the database
func (r *repository) Delete(key []byte) (err error) {
	bolt, err := r.ensureDataSource()
	if err != nil {
		return err
	}
	return bolt.Delete([]byte(r.Name()), key)
}



