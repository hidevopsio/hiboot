package db

import (
	"fmt"
)

const (
	DataSourceTypeBolt = "bolt"
	DataSourceTypeMysql = "mysql"
)

type DataSourceFactory struct {

}

func (dsf *DataSourceFactory) NewRepository(dataSourceType string, name string) (Repository, error)  {
	//switch dataSourceType {
	//case DataSourceTypeBolt:
	//	return , nil
	//}
	return nil, fmt.Errorf("database is not implemented")
}