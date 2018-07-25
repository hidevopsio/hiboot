package bolt

import (
	"github.com/hidevopsio/hiboot/pkg/starter"
	"github.com/hidevopsio/hiboot/pkg/starter/data"
)

type configuration struct {
	// the properties member name must be Bolt if the mapstructure is bolt,
	// so that the reference can be parsed
	Bolt properties `mapstructure:"bolt"`
}

func init() {
	starter.Add("bolt", configuration{})
}

func (c *configuration) dataSource() *bolt {
	bolt := GetInstance()
	bolt.Open(&c.Bolt)

	return bolt
}

func (c *configuration) NewRepository(name string) data.Repository {
	repo := new(repository)
	repo.SetDataSource(c.dataSource())
	repo.SetName(name)
	return repo
}
