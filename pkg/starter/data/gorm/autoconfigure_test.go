package gorm

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestConfiguration(t *testing.T) {
	configuration := new(configuration)
	configuration.Gorm = properties{
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

	repo := configuration.GormRepository()
	assert.NotEqual(t, nil, ds)
	err := repo.CloseDataSource()
	assert.Equal(t, nil, err)
}
