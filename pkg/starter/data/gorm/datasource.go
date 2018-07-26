package gorm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"fmt"
	"errors"
	"sync"
)

type DataSource interface {
	Open(p *properties) error
	IsOpened() bool
	Close() error
	DB() *gorm.DB
}

type dataSource struct {
	db *gorm.DB
}

var ds *dataSource
var dsOnce sync.Once

func GetDataSource() DataSource {
	dsOnce.Do(func() {
		ds = new(dataSource)
	})
	return ds
}

func (d *dataSource) Open(p *properties) error {
	var err error
	source := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=%v&loc=%v",
		p.Username, p.Password, p.Host, p.Port,  p.Database, p.Charset, p.ParseTime, p.Loc)

	d.db, err = gorm.Open(p.Type, source)

	if err != nil {
		defer d.db.Close()
		return err
	}

	return nil
}

func (d *dataSource) IsOpened() bool {
	return d.db != nil
}

func (d *dataSource) Close() error {
	if d.db != nil {
		err := d.db.Close()
		d.db = nil
		return err
	}
	return errors.New("database is not opened")
}

func (d *dataSource) DB() *gorm.DB {
	return d.db
}

