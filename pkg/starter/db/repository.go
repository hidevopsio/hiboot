package db

type Repository interface {
	SetName(name string)
	Name() string
	SetDataSource(dataSource interface{})
	DataSource() interface{}
}
