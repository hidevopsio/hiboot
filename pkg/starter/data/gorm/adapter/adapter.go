package adapter

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"errors"
)

type Gorm interface {
	Open(dialect string, args ...interface{}) (db *gorm.DB, err error)
	Close() error
}

type DB interface {
	AddError(err error) error
	AddForeignKey(field string, dest string, onDelete string, onUpdate string) *gorm.DB
	AddIndex(indexName string, columns ...string) *gorm.DB
	AddUniqueIndex(indexName string, columns ...string) *gorm.DB
	Assign(attrs ...interface{}) *gorm.DB
	Association(column string) *gorm.Association
	Attrs(attrs ...interface{}) *gorm.DB
	AutoMigrate(values ...interface{}) *gorm.DB
	Begin() *gorm.DB
	BlockGlobalUpdate(enable bool) *gorm.DB
	Callback() *gorm.Callback
	Close() error
	Commit() *gorm.DB
	CommonDB() gorm.SQLCommon
	Count(value interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	CreateTable(models ...interface{}) *gorm.DB
	DB() *sql.DB
	Debug() *gorm.DB
	Delete(value interface{}, where ...interface{}) *gorm.DB
	Dialect() gorm.Dialect
	DropColumn(column string) *gorm.DB
	DropTable(values ...interface{}) *gorm.DB
	DropTableIfExists(values ...interface{}) *gorm.DB
	Exec(sql string, values ...interface{}) *gorm.DB
	Find(out interface{}, where ...interface{}) *gorm.DB
	First(out interface{}, where ...interface{}) *gorm.DB
	FirstOrCreate(out interface{}, where ...interface{}) *gorm.DB
	FirstOrInit(out interface{}, where ...interface{}) *gorm.DB
	Get(name string) (value interface{}, ok bool)
	GetErrors() []error
	Group(query string) *gorm.DB
	HasBlockGlobalUpdate() bool
	HasTable(value interface{}) bool
	Having(query interface{}, values ...interface{}) *gorm.DB
	InstantSet(name string, value interface{}) *gorm.DB
	Joins(query string, args ...interface{}) *gorm.DB
	Last(out interface{}, where ...interface{}) *gorm.DB
	Limit(limit interface{}) *gorm.DB
	LogMode(enable bool) *gorm.DB
	Model(value interface{}) *gorm.DB
	ModifyColumn(column string, typ string) *gorm.DB
	New() *gorm.DB
	NewRecord(value interface{}) bool
	NewScope(value interface{}) *gorm.Scope
	Not(query interface{}, args ...interface{}) *gorm.DB
	Offset(offset interface{}) *gorm.DB
	Omit(columns ...string) *gorm.DB
	Or(query interface{}, args ...interface{}) *gorm.DB
	Order(value interface{}, reorder ...bool) *gorm.DB
	Pluck(column string, value interface{}) *gorm.DB
	Preload(column string, conditions ...interface{}) *gorm.DB
	//QueryExpr() *gorm.xpr
	Raw(sql string, values ...interface{}) *gorm.DB
	RecordNotFound() bool
	Related(value interface{}, foreignKeys ...string) *gorm.DB
	RemoveForeignKey(field string, dest string) *gorm.DB
	RemoveIndex(indexName string) *gorm.DB
	Rollback() *gorm.DB
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	Save(value interface{}) *gorm.DB
	Scan(dest interface{}) *gorm.DB
	ScanRows(rows *sql.Rows, result interface{}) error
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB
	Select(query interface{}, args ...interface{}) *gorm.DB
	Set(name string, value interface{}) *gorm.DB
	SetJoinTableHandler(source interface{}, column string, handler gorm.JoinTableHandlerInterface)
	//SetLogger(log logger)
	SingularTable(enable bool)
	//SubQuery() *expr
	Table(name string) *gorm.DB
	Take(out interface{}, where ...interface{}) *gorm.DB
	Unscoped() *gorm.DB
	Update(attrs ...interface{}) *gorm.DB
	UpdateColumn(attrs ...interface{}) *gorm.DB
	UpdateColumns(values interface{}) *gorm.DB
	Updates(values interface{}, ignoreProtectedAttrs ...bool) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
}

type GormDataSource struct {
	db *gorm.DB
}

var DatabaseIsNotOpenedError = errors.New("database is not opened")

func (d *GormDataSource) Open(dialect string, args ...interface{}) (db *gorm.DB, err error) {
	d.db, err = gorm.Open(dialect, args...)
	return d.db, err
}

func (d *GormDataSource) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return DatabaseIsNotOpenedError
}
