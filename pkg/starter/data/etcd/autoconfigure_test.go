package etcd

import (
	"go.etcd.io/etcd/clientv3"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/data/etcd/fake"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEtcd(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	conf := new(etcdConfiguration)

	t.Run("should create instance named etcdClient", func(t *testing.T) {
		conf.Properties = properties{
			Type:           "etcd",
			DialTimeout:    5,
			RequestTimeout: 10,
			Endpoints:      []string{"172.16.10.470:2379"},
			Cert: cert{CertFile: "config/certs/etcd.pem",
				KeyFile:       "config/certs/etcd-key.pem",
				TrustedCAFile: "config/certs/ca.pem"},
		}
		client := conf.Clientv3Client()
		assert.Equal(t, (*clientv3.Client)(nil), client)

	})

	t.Run("should not create instance named etcdRepository", func(t *testing.T) {
		client := new(clientv3.Client)
		repo := conf.EtcdRepository(client)
		assert.Equal(t, nil, repo)
	})

	t.Run("should create instance named etcdRepository", func(t *testing.T) {
		client := new(clientv3.Client)
		client.KV = new(fake.Repository)
		repo := conf.EtcdRepository(client)
		assert.Equal(t, client.KV, repo)
	})
}
