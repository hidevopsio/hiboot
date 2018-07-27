package gorm

import (
	"github.com/hidevopsio/hiboot/pkg/starter"
	"github.com/hidevopsio/hiboot/pkg/starter/data"
)

type configuration struct {
	// the properties member name must be Gorm if the mapstructure is gorm,
	// so that the reference can be parsed
	GormProperties properties `mapstructure:"gorm"`
}

func init() {
	starter.Add("gorm", configuration{})
}

func (c *configuration) dataSource() DataSource {
	dataSource := GetDataSource()
	if ! dataSource.IsOpened() {
		dataSource.Open(&c.GormProperties)
	}
	return dataSource
}

// GormRepository method name must be unique
func (c *configuration) GormRepository() data.Repository {
	repository := GetRepository()
	repository.SetDataSource(c.dataSource())
	return repository
}
