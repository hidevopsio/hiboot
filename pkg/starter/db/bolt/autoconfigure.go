package bolt

import (
	"github.com/hidevopsio/hiboot/pkg/starter"
)

type Configuration struct {
	BoltProperties Properties `mapstructure:"bolt"`
}

func init() {
	starter.Add("bolt", Configuration{})
}

func (c *Configuration) dataSource() *Bolt {
	bolt := GetInstance()
	bolt.Open(&c.BoltProperties)

	return bolt
}

func (c *Configuration) NewRepository(name string) RepositoryInterface {
	repo := new(Repository)
	repo.SetDataSource(c.dataSource())
	repo.SetName(name)
	return repo
}
