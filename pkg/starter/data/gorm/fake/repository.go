// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fake

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type Repository struct {
}

// New clone a new db connection without search conditions
func (s *Repository) New() *gorm.DB {
	return nil
}

// Close close current db connection.  If database connection is not an io.Closer, returns an error.
func (s *Repository) Close() error {
	return nil
}

// DB get `*sql.DB` from current connection
// If the underlying database connection is not a *sql.DB, returns nil
func (s *Repository) DB() *sql.DB {
	return nil
}

// CommonDB return the underlying `*sql.DB` or `*sql.Tx` instance, mainly intended to allow coexistence with legacy non-GORM code.
func (s *Repository) CommonDB() gorm.SQLCommon {
	return nil
}

// Dialect get dialect
func (s *Repository) Dialect() gorm.Dialect {
	return nil
}

func (s *Repository) Callback() *gorm.Callback {
	return nil
}

// LogMode set log mode, `true` for detailed logs, `false` for no log, default, will only print error logs
func (s *Repository) LogMode(enable bool) *gorm.DB {
	return nil
}

// BlockGlobalUpdate if true, generates an error on update/delete without where clause.
// This is to prevent eventual error with empty objects updates/deletions
func (s *Repository) BlockGlobalUpdate(enable bool) *gorm.DB {
	return nil
}

// HasBlockGlobalUpdate return state of block
func (s *Repository) HasBlockGlobalUpdate() bool {
	return false
}

// SingularTable use singular table by default
func (s *Repository) SingularTable(enable bool) {
}

// NewScope create a scope for current operation
func (s *Repository) NewScope(value interface{}) *gorm.Scope {
	return nil
}

// Where return a new relation, filter records with given conditions, accepts `map`, `struct` or `string` as conditions, refer http://jinzhu.github.io/gorm/crud.html#query
func (s *Repository) Where(query interface{}, args ...interface{}) *gorm.DB {
	return nil
}

// Or filter records that match before conditions or this one, similar to `Where`
func (s *Repository) Or(query interface{}, args ...interface{}) *gorm.DB {
	return nil
}

// Not filter records that don't match current conditions, similar to `Where`
func (s *Repository) Not(query interface{}, args ...interface{}) *gorm.DB {
	return nil
}

// Limit specify the number of records to be retrieved
func (s *Repository) Limit(limit interface{}) *gorm.DB {
	return nil
}

// Offset specify the number of records to skip before starting to return the records
func (s *Repository) Offset(offset interface{}) *gorm.DB {
	return nil
}

// Order specify order when retrieve records from database, set reorder to `true` to overwrite defined conditions
//     db.Order("name DESC")
//     db.Order("name DESC", true) // reorder
//     db.Order(gorm.Expr("name = ? DESC", "first")) // sql expression
func (s *Repository) Order(value interface{}, reorder ...bool) *gorm.DB {
	return nil
}

// Select specify fields that you want to retrieve from database when querying, by default, will select all fields;
// When creating/updating, specify fields that you want to save to database
func (s *Repository) Select(query interface{}, args ...interface{}) *gorm.DB {
	return nil
}

// Omit specify fields that you want to ignore when saving to database for creating, updating
func (s *Repository) Omit(columns ...string) *gorm.DB {
	return nil
}

// Group specify the group method on the find
func (s *Repository) Group(query string) *gorm.DB {
	return nil
}

// Having specify HAVING conditions for GROUP BY
func (s *Repository) Having(query interface{}, values ...interface{}) *gorm.DB {
	return nil
}

// Joins specify Joins conditions
//     db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org").Find(&user)
func (s *Repository) Joins(query string, args ...interface{}) *gorm.DB {
	return nil
}

func (s *Repository) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	return nil
}

// Unscoped return all record including deleted record, refer Soft Delete https://jinzhu.github.io/gorm/crud.html#soft-delete
func (s *Repository) Unscoped() *gorm.DB {
	return nil
}

// Attrs initialize struct with argument if record not found with `FirstOrInit` https://jinzhu.github.io/gorm/crud.html#firstorinit or `FirstOrCreate` https://jinzhu.github.io/gorm/crud.html#firstorcreate
func (s *Repository) Attrs(attrs ...interface{}) *gorm.DB {
	return nil
}

// Assign assign result with argument regardless it is found or not with `FirstOrInit` https://jinzhu.github.io/gorm/crud.html#firstorinit or `FirstOrCreate` https://jinzhu.github.io/gorm/crud.html#firstorcreate
func (s *Repository) Assign(attrs ...interface{}) *gorm.DB {
	return nil
}

// First find first record that match given conditions, order by primary key
func (s *Repository) First(out interface{}, where ...interface{}) *gorm.DB {
	return nil
}

// Take return a record that match given conditions, the order will depend on the database implementation
func (s *Repository) Take(out interface{}, where ...interface{}) *gorm.DB {
	return nil
}

// Last find last record that match given conditions, order by primary key
func (s *Repository) Last(out interface{}, where ...interface{}) *gorm.DB {
	return nil
}

// Find find records that match given conditions
func (s *Repository) Find(out interface{}, where ...interface{}) *gorm.DB {
	return nil
}

// Scan scan value to a struct
func (s *Repository) Scan(dest interface{}) *gorm.DB {
	return nil
}

// Row return `*sql.Row` with given conditions
func (s *Repository) Row() *sql.Row {
	return nil
}

// Rows return `*sql.Rows` with given conditions
func (s *Repository) Rows() (*sql.Rows, error) {
	return nil, nil
}

// ScanRows scan `*sql.Rows` to give struct
func (s *Repository) ScanRows(rows *sql.Rows, result interface{}) error {
	return nil
}

// Pluck used to query single column from a model as a map
//     var ages []int64
//     db.Find(&users).Pluck("age", &ages)
func (s *Repository) Pluck(column string, value interface{}) *gorm.DB {
	return nil
}

// Count get how many records for a model
func (s *Repository) Count(value interface{}) *gorm.DB {
	return nil
}

// Related get related associations
func (s *Repository) Related(value interface{}, foreignKeys ...string) *gorm.DB {
	return nil
}

// FirstOrInit find first matched record or initialize a new one with given conditions (only works with struct, map conditions)
// https://jinzhu.github.io/gorm/crud.html#firstorinit
func (s *Repository) FirstOrInit(out interface{}, where ...interface{}) *gorm.DB {
	return nil
}

// FirstOrCreate find first matched record or create a new one with given conditions (only works with struct, map conditions)
// https://jinzhu.github.io/gorm/crud.html#firstorcreate
func (s *Repository) FirstOrCreate(out interface{}, where ...interface{}) *gorm.DB {
	return nil
}


// Update update attributes with callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
func (s *Repository) Update(attrs ...interface{}) *gorm.DB {
	return nil
}

// Updates update attributes with callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
func (s *Repository) Updates(values interface{}, ignoreProtectedAttrs ...bool) *gorm.DB {
	return nil
}

// UpdateColumn update attributes without callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
func (s *Repository) UpdateColumn(attrs ...interface{}) *gorm.DB {
	return nil
}

// UpdateColumns update attributes without callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
func (s *Repository) UpdateColumns(values interface{}) *gorm.DB {
	return nil
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (s *Repository) Save(value interface{}) *gorm.DB {
	return nil
}

// Create insert the value into database
func (s *Repository) Create(value interface{}) *gorm.DB {
	return nil
}


// Delete delete value match given conditions, if the value has primary key, then will including the primary key as condition
func (s *Repository) Delete(value interface{}, where ...interface{}) *gorm.DB {
	return nil
}

// Raw use raw sql as conditions, won't run it unless invoked by other methods
//    db.Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result)
func (s *Repository) Raw(sql string, values ...interface{}) *gorm.DB {
	return nil
}

// Exec execute raw sql
func (s *Repository) Exec(sql string, values ...interface{}) *gorm.DB {
	return nil
}

// Model specify the model you would like to run db operations
//    // update all users's name to `hello`
//    db.Model(&User{}).Update("name", "hello")
//    // if user's primary key is non-blank, will use it as condition, then will only update the user's name to `hello`
//    db.Model(&user).Update("name", "hello")
func (s *Repository) Model(value interface{}) *gorm.DB {
	return nil
}

// Table specify the table you would like to run db operations
func (s *Repository) Table(name string) *gorm.DB {
	return nil
}


// Debug start debug mode
func (s *Repository) Debug() *gorm.DB {
	return nil
}

// Begin begin a transaction
func (s *Repository) Begin() *gorm.DB {
	return nil
}

// Commit commit a transaction
func (s *Repository) Commit() *gorm.DB {
	return nil
}

// Rollback rollback a transaction
func (s *Repository) Rollback() *gorm.DB {
	return nil
}

// NewRecord check if value's primary key is blank
func (s *Repository) NewRecord(value interface{}) bool {
	return s.NewScope(value).PrimaryKeyZero()
}

// RecordNotFound check if returning ErrRecordNotFound error
func (s *Repository) RecordNotFound() bool {
	return false
}

// CreateTable create table for models
func (s *Repository) CreateTable(models ...interface{}) *gorm.DB {
	return nil
}

// DropTable drop table for models
func (s *Repository) DropTable(values ...interface{}) *gorm.DB {
	return nil
}

// DropTableIfExists drop table if it is exist
func (s *Repository) DropTableIfExists(values ...interface{}) *gorm.DB {
	return nil
}


// HasTable check has table or not
func (s *Repository) HasTable(value interface{}) bool {
	return false
}

// AutoMigrate run auto migration for given models, will only add missing fields, won't delete/change current data
func (s *Repository) AutoMigrate(values ...interface{}) *gorm.DB {
	return nil
}

// ModifyColumn modify column to type
func (s *Repository) ModifyColumn(column string, typ string) *gorm.DB {
	return nil
}

// DropColumn drop a column
func (s *Repository) DropColumn(column string) *gorm.DB {
	return nil
}

// AddIndex add index for columns with given name
func (s *Repository) AddIndex(indexName string, columns ...string) *gorm.DB {
	return nil
}

// AddUniqueIndex add unique index for columns with given name
func (s *Repository) AddUniqueIndex(indexName string, columns ...string) *gorm.DB {
	return nil
}

// RemoveIndex remove index with name
func (s *Repository) RemoveIndex(indexName string) *gorm.DB {
	return nil
}

// AddForeignKey Add foreign key to the given scope, e.g:
//     db.Model(&User{}).AddForeignKey("city_id", "cities(id)", "RESTRICT", "RESTRICT")
func (s *Repository) AddForeignKey(field string, dest string, onDelete string, onUpdate string) *gorm.DB {
	return nil
}

// RemoveForeignKey Remove foreign key from the given scope, e.g:
//     db.Model(&User{}).RemoveForeignKey("city_id", "cities(id)")
func (s *Repository) RemoveForeignKey(field string, dest string) *gorm.DB {
	return nil
}

// Association start `Association Mode` to handler relations things easir in that mode, refer: https://jinzhu.github.io/gorm/associations.html#association-mode
func (s *Repository) Association(column string) *gorm.Association {
	return nil
}

// Preload preload associations with given conditions
//    db.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
func (s *Repository) Preload(column string, conditions ...interface{}) *gorm.DB {
	return nil
}

// Set set setting by name, which could be used in callbacks, will clone a new db, and update its setting
func (s *Repository) Set(name string, value interface{}) *gorm.DB {
	return nil
}

// InstantSet instant set setting, will affect current db
func (s *Repository) InstantSet(name string, value interface{}) *gorm.DB {
	return nil
}

// Get get setting by name
func (s *Repository) Get(name string) (value interface{}, ok bool) {
	return nil, false
}

// SetJoinTableHandler set a model's join table handler for a relation
func (s *Repository) SetJoinTableHandler(source interface{}, column string, handler gorm.JoinTableHandlerInterface) {
}

// AddError add error to the db
func (s *Repository) AddError(err error) error {
	return nil
}

// GetErrors get happened errors from the db
func (s *Repository) GetErrors() []error {
	return nil
}
