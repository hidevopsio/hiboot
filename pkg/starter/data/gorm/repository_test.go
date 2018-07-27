package gorm

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/starter/data/gorm/adapter"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/starter/data"
	"os"
)

func TestRepository(t *testing.T) {
	gorm := &properties{
		Type:      "mysql",
		Host:      "mysql-dev",
		Port:      "3306",
		Username:  os.Getenv("MYSQL_USERNAME"),
		Password:  os.Getenv("MYSQL_PASSWORD"),
		Database:  "test",
		ParseTime: "True",
		Charset:   "utf8",
		Loc:       "Asia%2FShanghai",
	}
	dataSource := new(dataSource)
	dataSource.gorm = new(adapter.FakeDataSource)
	repo := new(repository)

	t.Run("should report error if trying to close the unopened database", func(t *testing.T) {
		err := repo.CloseDataSource()
		assert.Equal(t, data.InvalidDataSourceError, err)
	})

	dataSource.Open(gorm)
	repo.SetDataSource(dataSource)

	t.Run("should close database properly", func(t *testing.T) {
		err := repo.CloseDataSource()
		assert.Equal(t, nil, err)
	})
}