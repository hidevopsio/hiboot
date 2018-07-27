package adapter

import "github.com/jinzhu/gorm"

type FakeDataSource struct {
	db *gorm.DB
}

func (d *FakeDataSource) Open(dialect string, args ...interface{}) (db *gorm.DB, err error) {
	return nil, nil
}

func (d *FakeDataSource) Close() error {
	return nil
}

