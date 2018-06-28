package db

type KVRepository interface {
	Put(bucket, key, value []byte) error
	Get(bucket, key []byte) ([]byte, error)
	Delete(bucket, key []byte) error
}
