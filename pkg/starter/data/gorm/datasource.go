package gorm

import (
	"fmt"
	"errors"
	"sync"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/hidevopsio/hiboot/pkg/starter/data/gorm/adapter"
)

type DataSource interface {
	Open(p *properties) error
	IsOpened() bool
	Close() error
	DB() adapter.DB
}

type dataSource struct {
	gorm adapter.Gorm
	db adapter.DB
}

var DatabaseIsNotOpenedError = errors.New("database is not opened")
var ds *dataSource
var dsOnce sync.Once

func GetDataSource() DataSource {
	dsOnce.Do(func() {
		ds = new(dataSource)
		ds.gorm = new(adapter.GormDataSource)
	})
	return ds
}

func (d *dataSource) Open(p *properties) error {
	var err error
	source := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=%v&loc=%v",
		p.Username, p.Password, p.Host, p.Port,  p.Database, p.Charset, p.ParseTime, p.Loc)

	d.db, err = d.gorm.Open(p.Type, source)

	if err != nil {
		defer d.gorm.Close()
		return err
	}

	return nil
}

func (d *dataSource) IsOpened() bool {
	return d.db != nil
}

func (d *dataSource) Close() error {
	if d.db != nil {
		err := d.gorm.Close()
		d.db = nil
		return err
	}
	return DatabaseIsNotOpenedError
}

func (d *dataSource) DB() adapter.DB {
	return d.db
}

