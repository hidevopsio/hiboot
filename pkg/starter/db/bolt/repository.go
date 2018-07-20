package bolt

import "errors"

type RepositoryInterface interface {
	Put(key, value []byte) error
	Get(key []byte) (result []byte, err error)
	Delete(key []byte) (err error)
}

type Repository struct {
	dataSource *Bolt
	name []byte
}

var NilPointError = errors.New("dataSource is nil")

// SetName set repository name
func (r *Repository) SetName(name string)  {
	r.name = []byte(name)
}

// SetDataSource set repository data source
func (r *Repository) SetDataSource(dataSource interface{}) {
	r.dataSource = dataSource.(*Bolt)
}

// Name get the repository name
func (r *Repository) Name() string  {
	return string(r.name)
}

// DataSource get the repository data source
func (r *Repository) DataSource() interface{}  {
	return r.dataSource
}

// Put inserts a key:value pair into the database
func (r *Repository) Put(key, value []byte) error {
	if r.dataSource == nil {
		return NilPointError
	}
	return r.dataSource.Put(r.name, key, value)
}

// Get retrieves a key:value pair from the database
func (r *Repository) Get(key []byte) (result []byte, err error)  {
	if r.dataSource == nil {
		return nil, NilPointError
	}
	return r.dataSource.Get(r.name, key)
}

// Delete removes a key:value pair from the database
func (r *Repository) Delete(key []byte) (err error) {
	if r.dataSource == nil {
		return NilPointError
	}
	return r.dataSource.Delete(r.name, key)
}

// Close close the data source
func (r *Repository) Close() (err error) {
	if r.dataSource == nil {
		return NilPointError
	}
	return r.dataSource.Close()
}


