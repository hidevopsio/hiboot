package factory

import (
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/gotest"
)

type InstanceFactory struct {
	instances cmap.ConcurrentMap
}

func (f *InstanceFactory) Init(instances cmap.ConcurrentMap)  {
	f.instances = instances
}

func (f *InstanceFactory) SetInstance(name string, instance interface{}) {
	if _, ok := f.instances.Get(name); ok && !gotest.IsRunning() {
		log.Fatalf("[factory] instance name % is already taken", name)
	}
	f.instances.Set(name, instance)
}

func (f *InstanceFactory) GetInstance(name string) (inst interface{}) {
	var ok bool
	if inst, ok = f.instances.Get(name); !ok {
		return nil
	}
	return
}