package bolt

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/starter/data"
)


type User struct {
	ID string
	Name string
	Age uint
}

type Foo struct {
	Name string
}
type Bar struct {
	// TODO: find model ID as the key of current record
	Name string	`model:"ID"`
}

var (
	id = "jd"
	user = &User{ID: id, Name: "John Doe", Age: 18}
)

func TestRepositoryCrd(t *testing.T) {

	properties := &properties{
		Database: "test.db",
		Mode: 0600,
		Timeout: 2,
	}

	r := GetRepository()
	d := GetDataSource()


	t.Run("should open bolt database", func(t *testing.T) {
		err := d.Open(nil)
		assert.Equal(t, InvalidPropertiesError, err)
	})

	t.Run("should open bolt database", func(t *testing.T) {
		err := d.Open(properties)
		assert.Equal(t, nil, err)
	})

	r.SetDataSource(d)

	t.Run("should put data into bolt database", func(t *testing.T) {
		err := r.Put(user)
		assert.Equal(t, nil, err)
	})

	t.Run("should get data into bolt database", func(t *testing.T) {
		u := &User{ID: id}
		err := r.Get(u)
		assert.Equal(t, nil, err)
		assert.Equal(t, u.Name, user.Name)
	})

	t.Run("should delete data into bolt database", func(t *testing.T) {
		u := &User{ID: id}
		err := r.Delete(u)
		assert.Equal(t, nil, err)
	})

	t.Run("should put data into bolt database with key", func(t *testing.T) {
		err := r.Put("newKey", user)
		assert.Equal(t, nil, err)

		u := &User{}
		err = r.Get("newKey", u)
		assert.Equal(t, nil, err)
		assert.Equal(t, u.Name, user.Name)
	})

	t.Run("should return InvalidDataModelError", func(t *testing.T) {
		err := r.Put(&Foo{Name: "foo"})
		assert.Equal(t, data.InvalidDataModelError, err)
	})

	// close bolt database
	r.CloseDataSource()
}

func TestRepositoryWithNilDataSource(t *testing.T) {
	r := &repository{}

	t.Run("should put data into bolt database", func(t *testing.T) {
		err := r.Put(user)
		assert.Equal(t, data.InvalidDataSourceError, err)
	})

	t.Run("should get data into bolt database", func(t *testing.T) {
		u := &User{ID: id}
		err := r.Get(u)
		assert.Equal(t, data.InvalidDataSourceError, err)
	})

	t.Run("should delete data into bolt database", func(t *testing.T) {
		u := &User{ID: id}
		err := r.Delete(u)
		assert.Equal(t, data.InvalidDataSourceError, err)
	})
}
