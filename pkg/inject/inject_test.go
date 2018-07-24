package inject

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/starter/db/bolt"
	"github.com/hidevopsio/hiboot/pkg/starter/db"
	"github.com/hidevopsio/hiboot/pkg/starter"
)

type User struct {
	Name string
}

type fakeRepository struct {
	db.BaseRepository
}

// Put inserts a key:value pair into the database
func (r *fakeRepository) Put(key, value []byte) error {
	return nil
}

// Get retrieves a key:value pair from the database
func (r *fakeRepository) Get(key []byte) (result []byte, err error)  {
	return []byte("fake data"), nil
}

// Delete removes a key:value pair from the database
func (r *fakeRepository) Delete(key []byte) (err error) {
	return nil
}

type fakeConfiguration struct{
	
}

func (c *fakeConfiguration) NewRepository(name string) db.Repository {
	repo := new(fakeRepository)
	repo.SetName(name)
	return repo
}


type UserService struct {
	User           *User           `inject:"user"`
	UserRepository bolt.Repository `inject:"userRepository,dataSourceType=fake,namespace=user"`
	Url            string          `value:"${fake.url:http://localhost:8080}"`
}

type RecursiveInjectTest struct {
	UserService *UserService
}

func init() {
	starter.Add("fake", fakeConfiguration{})
	starter.GetAutoConfiguration().Build()
}

func TestNotInject(t *testing.T) {
	baz := new(UserService)
	assert.Equal(t, (*User)(nil), baz.User)
}

func TestInject(t *testing.T) {
	t.Run("should inject repository", func(t *testing.T) {
		us := new(UserService)
		IntoObject(reflect.ValueOf(us))
		assert.NotEqual(t, (*User)(nil), us.User)
		assert.NotEqual(t, (*fakeRepository)(nil), us.UserRepository)
	})

	t.Run("should inject recursively", func(t *testing.T) {
		ps := &RecursiveInjectTest{UserService: new(UserService)}
		IntoObject(reflect.ValueOf(ps))
		assert.NotEqual(t, (*User)(nil), ps.UserService.User)
		assert.NotEqual(t, (*fakeRepository)(nil), ps.UserService.UserRepository)
	})
}
