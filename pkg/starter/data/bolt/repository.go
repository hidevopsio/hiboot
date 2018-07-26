package bolt

import (
	"errors"
	"sync"
	"github.com/boltdb/bolt"
	"github.com/hidevopsio/hiboot/pkg/starter/data"
	"encoding/json"
)

type Repository interface {
	data.KVRepository
}

type repository struct {
	data.BaseKVRepository
	db *bolt.DB
	dataSource DataSource
}

var repo *repository
var once sync.Once
var InvalidPropertiesError = errors.New("properties must not be nil")

func GetRepository() *repository {
	once.Do(func() {
		repo = &repository{}
	})
	return repo
}

func (r *repository) parse(params []interface{}) ([]byte, []byte, interface{}, error)  {
	if r.db == nil {
		return nil, nil, nil, data.InvalidDataSourceError
	}
	return r.Parse(params)
}

// Open bolt database
func (r *repository) SetDataSource(d interface{})  {
	if d != nil {
		r.dataSource = d.(DataSource)
		r.db = r.dataSource.DB()
	}
}

func (r *repository) DataSource() interface{}  {
	return r.dataSource
}

func (r *repository) CloseDataSource() error {
	if r.dataSource != nil {
		return r.dataSource.Close()
	}
	return data.InvalidDataSourceError
}

// Put inserts a key:value pair into the database
func (r *repository) Put(params ...interface{}) error {
	bucketName, key, value, err := r.parse(params)
	if err != nil {
		return err
	}

	err = r.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		// marshal data to bytes
		b, err := json.Marshal(value)

		err = bucket.Put(key, b)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// Get retrieves a key:value pair from the database
func (r *repository) Get(params ...interface{}) error {
	bucketName, key, value, err := r.parse(params)
	if err != nil {
		return err
	}
	var result []byte
	err = r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b != nil {
			v := b.Get(key)
			if v != nil {
				result = make([]byte, len(v))
				copy(result, v)
			}
		} else {
			result = []byte("")
		}
		return nil
	})
	if err == nil {
		err = json.Unmarshal(result, value)
	}
	return err
}

// Delete removes a key:value pair from the database
func (r *repository) Delete(params ...interface{}) error {
	bucketName, key, _, err := r.parse(params)
	if err != nil {
		return err
	}
	err = r.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		err = bucket.Delete(key)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}


