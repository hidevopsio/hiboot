package data

// KVRepository is the Key/Value Repository interface
type KVRepository interface {
	Repository
	Put(key, value []byte) error
	Get(key []byte) (result []byte, err error)
	Delete(key []byte) (err error)
}

type BaseKVRepository struct {
	BaseRepository
}

// Put inserts a key:value pair into the database
func (r *BaseKVRepository) Put(key, value []byte) error {
	return nil
}

// Get retrieves a key:value pair from the database
func (r *BaseKVRepository) Get(key []byte) (result []byte, err error)  {
	return []byte(""), nil
}

// Delete removes a key:value pair from the database
func (r *BaseKVRepository) Delete(key []byte) (err error) {
	return nil
}