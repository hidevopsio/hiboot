package db

type DataSourceInterface interface {
	Open(dataSource interface{}) error
	Close() error
	SetNamespace(name string)
}
