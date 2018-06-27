package db

import (
	"github.com/hidevopsio/hiboot/pkg/starter/db/bolt"
	"fmt"
)

const (
	DataSourceTypeBolt = "bolt"
	DataSourceTypeMysql = "mysql"
)

type DataSourceFactory struct {

}

func (dsf *DataSourceFactory) New(dataSourceType string) (DataSourceInterface, error)  {
	switch dataSourceType {
	case DataSourceTypeBolt:
		return new(bolt.Bolt), nil
	}
	return nil, fmt.Errorf("database is not implemented")
}