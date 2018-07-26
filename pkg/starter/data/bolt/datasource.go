package bolt

import (
	"github.com/boltdb/bolt"
	"time"
	"sync"
)

type DataSource interface {
	Open(properties *properties) error
	Close() error
	IsOpened() bool
	DB() *bolt.DB
}

type dataSource struct {
	db *bolt.DB
}

var ds *dataSource
var do sync.Once

func GetDataSource() *dataSource {
	do.Do(func() {
		ds = new(dataSource)
	})
	return ds
}

func (d *dataSource) DB() *bolt.DB {
	return d.db
}

func (d *dataSource) Open(properties *properties) error {
	if properties == nil {
		return InvalidPropertiesError
	}

	var err error
	d.db, err = bolt.Open(properties.Database,
		properties.Mode,
		&bolt.Options{Timeout: time.Duration(properties.Timeout) * time.Second},
	)

	if err != nil {
		if d.db != nil {
			defer d.db.Close()
		}
	}

	return err
}

// Close database
func (d *dataSource) Close() error {
	err := d.db.Close()
	d.db = nil
	return err
}

// IsOpened check if db is opened
func (d *dataSource) IsOpened() bool {
	return d.db != nil
}