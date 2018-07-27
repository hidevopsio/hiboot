package gorm

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/data/gorm/adapter"
)

type User struct {
	ID   uint `gorm:"primary_key"`
	Username string
	Password string
	Age  int `gorm:"default:18"`
	Gender  int `gorm:"default:18"`
}

func (User) TableName() string {
	return "user"
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestDataSourceOpen(t *testing.T) {
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

	t.Run("should report error when close database if it's not opened", func(t *testing.T) {
		err := dataSource.Close()
		assert.Equal(t, DatabaseIsNotOpenedError, err)
	})

	t.Run("should open database", func(t *testing.T) {
		err := dataSource.Open(gorm)
		assert.Equal(t, nil, err)
	})

	t.Run("should close database", func(t *testing.T) {
		err := dataSource.Close()
		assert.Equal(t, nil, err)
	})

}
