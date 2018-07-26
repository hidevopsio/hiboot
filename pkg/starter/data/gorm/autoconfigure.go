package gorm

import (
	"github.com/hidevopsio/hiboot/pkg/starter"
	"github.com/hidevopsio/hiboot/pkg/starter/data"
)

type configuration struct {
	// the properties member name must be Gorm if the mapstructure is gorm,
	// so that the reference can be parsed
	Gorm properties `mapstructure:"gorm"`
}

func init() {
	starter.Add("gorm", configuration{})
}

func (c *configuration) dataSource() DataSource {
	ds := GetDataSource()
	if ! ds.IsOpened() {
		ds.Open(&c.Gorm)
	}
	return ds
}

func (c *configuration) GormRepository() data.Repository {
	repo := GetRepository()
	repo.SetDataSource(c.dataSource())
	return repo
}
