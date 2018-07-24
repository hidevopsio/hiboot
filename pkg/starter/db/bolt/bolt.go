package bolt

import (
	"time"
	"github.com/hidevopsio/hiboot/pkg/log"
	boltdb "github.com/boltdb/bolt"
	"errors"
	"sync"
)

type bolt struct {
	DB        *boltdb.DB
	BK        *boltdb.Bucket
	BS        *boltdb.BucketStats
}

var instance *bolt
var once sync.Once

func GetInstance() *bolt {
	once.Do(func() {
		instance = &bolt{}
	})
	return instance
}

// Open bolt database
func (b *bolt) Open(properties *properties) error {

	if b.DB != nil {
		return nil
	}

	if properties == nil {
		return errors.New("properties must not be nil")
	}

	var err error
	b.DB, err = boltdb.Open(properties.Database,
		properties.Mode,
		&boltdb.Options{Timeout: time.Duration(properties.Timeout) * time.Second},
	)

	if err != nil {
		defer b.DB.Close()
		log.Fatal(err)
	}

	return err
}

// Close database
func (b *bolt) Close() error {
	err := b.DB.Close()
	b.DB = nil
	return err
}

// Put inserts a key:value pair into the database
func (b *bolt) Put(bucket, key, value []byte) error {
	//dbPath := bt.db.Path()
	//log.Println("DB Info: ", reflect.TypeOf(dbPath), dbPath)
	err := b.DB.Update(func(tx *boltdb.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// Get retrieves a key:value pair from the database
func (b *bolt) Get(bucket, key []byte) (result []byte, err error)  {
	b.DB.View(func(tx *boltdb.Tx) error {
		b := tx.Bucket(bucket)
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
	return
}

// DeleteKey removes a key:value pair from the database
func (b *bolt) Delete(bucket, key []byte) (err error) {

	err = b.DB.Update(func(tx *boltdb.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket)
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
