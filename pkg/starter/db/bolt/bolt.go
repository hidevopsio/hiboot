package bolt

import (
	"time"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/mitchellh/mapstructure"
	boltdb "github.com/boltdb/bolt"
)

type Bolt struct {
	DB *boltdb.DB
	BK *boltdb.Bucket
	BS *boltdb.BucketStats
}

type DataSource struct {
	Database string        `json:"database"`
	Mode     os.FileMode   `json:"mode"`
	Timeout  int64 `json:"timeout"`
}

// Open bolt database
func (b *Bolt) Open(dataSource interface{}) error {

	var ds DataSource
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &ds,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	err = decoder.Decode(dataSource)
	if err != nil {
		return err
	}

	b.DB, err = boltdb.Open(ds.Database,
		ds.Mode,
		&boltdb.Options{Timeout: time.Duration(ds.Timeout) * time.Second},
	)

	if err != nil {
		defer b.DB.Close()
		log.Fatal(err)
	}

	return err
}

// Close database
func (b *Bolt) Close() error {
	return b.DB.Close()
}

// Put inserts a key:value pair into the database
func (b *Bolt) Put(bucket, key, value []byte) error {
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
func (b *Bolt) Get(bucket, key []byte) (result []byte, err error)  {
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
func (b *Bolt) Delete(bucket, key []byte) (err error) {

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
