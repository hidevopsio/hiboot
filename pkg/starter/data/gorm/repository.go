package gorm

import (
	"github.com/hidevopsio/hiboot/pkg/starter/data"
	"sync"
	"github.com/hidevopsio/hiboot/pkg/starter/data/gorm/adapter"
)

type Repository interface {
	data.Repository
}

type repository struct {
	data.BaseRepository
	dataSource DataSource
	db adapter.DB
}

var repo *repository
var once sync.Once

func GetRepository() *repository {
	once.Do(func() {
		repo = &repository{}
	})
	return repo
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