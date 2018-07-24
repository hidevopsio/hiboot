package db

type Repository interface {
	SetName(name string)
	Name() string
	SetDataSource(dataSource interface{})
	DataSource() interface{}
}

type BaseRepository struct {
	dataSource interface{}
	name string
}

// SetName set repository name
func (r *BaseRepository) SetName(name string)  {
	r.name = name
}

// SetDataSource set repository data source
func (r *BaseRepository) SetDataSource(dataSource interface{}) {
	r.dataSource = dataSource
}

// Name get the repository name
func (r *BaseRepository) Name() string  {
	return string(r.name)
}

// DataSource get the repository data source
func (r *BaseRepository) DataSource() interface{}  {
	return r.dataSource
}