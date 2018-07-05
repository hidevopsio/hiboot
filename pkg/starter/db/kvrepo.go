package db


// KVRepository is the Key/Value Repository interface
type KVRepository interface {
	// Put key value pair to specific bucket
	Put(key, value []byte) error
	// Get value from specific bucket with key
	Get(key []byte) ([]byte, error)
	// Delete key in specific bucket
	Delete(key []byte) error
}
