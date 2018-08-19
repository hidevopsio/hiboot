package fake

import "github.com/hidevopsio/gorm"

type DataSource struct {
}

func (d *DataSource) Open(dialect string, args ...interface{}) (db gorm.Repository, err error){
	return nil, nil
}

func (d *DataSource) Close() error {
	return nil
}


func (d *DataSource) IsOpened() bool {
	return false
}

func (d *DataSource) Repository() gorm.Repository {
	return nil
}