package data

import "errors"

var InvalidDataSourceError = errors.New("invalid dataSource")
var InvalidDataModelError = errors.New("invalid data model, must contains string field ID and assigns string value")
var NotImplemenedError = errors.New("method is not implemented")

type Repository interface {
	SetDataSource(dataSource interface{})
	DataSource() interface{}
	CloseDataSource() error
}

type BaseRepository struct {
}

func (r *BaseRepository) SetDataSource(dataSource interface{})  {

}

func (r *BaseRepository) DataSource() interface{}  {
	return NotImplemenedError
}

func (r *BaseRepository) CloseDataSource() error  {
	return NotImplemenedError
}