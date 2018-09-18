package fake

import (
	"github.com/stretchr/testify/mock"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
)

type Repository struct {
	mock.Mock
}

func (e *Repository) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	args := e.Called(nil, key)
	return args[0].(*clientv3.PutResponse), args.Error(1)
}

func (e *Repository) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	args := e.Called(nil, key)
	return args[0].(*clientv3.GetResponse), args.Error(1)
}

func (e *Repository) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	args := e.Called(nil, key)
	return args[0].(*clientv3.DeleteResponse), args.Error(1)
}

func (e *Repository) Compact(ctx context.Context, rev int64, opts ...clientv3.CompactOption) (*clientv3.CompactResponse, error) {
	return nil, nil
}

func (e *Repository) Do(ctx context.Context, op clientv3.Op) (clientv3.OpResponse, error) {
	return clientv3.OpResponse{}, nil
}

func (e *Repository) Txn(ctx context.Context) clientv3.Txn {
	return nil
}
