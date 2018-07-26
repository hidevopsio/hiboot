package bolt

import (
	"github.com/hidevopsio/hiboot/pkg/starter"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type configuration struct {
	// the properties member name must be Bolt if the mapstructure is bolt,
	// so that the reference can be parsed
	Bolt properties `mapstructure:"bolt"`
}

func init() {
	starter.Add("bolt", configuration{})
}

func (c *configuration) dataSource() DataSource {
	dataSource := GetDataSource()
	if ! dataSource.IsOpened() {
		err := dataSource.Open(&c.Bolt)
		if err != nil {
			log.Error(err.Error())
		}
	}
	return dataSource
}

func (c *configuration) BoltRepository() Repository {
	repository := GetRepository()
	repository.SetDataSource(c.dataSource())
	return repository
}
