package instance

import (
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/gotest"
)

type Factory interface {
	Initialized() bool
	SetInstance(name string, instance interface{})
	GetInstance(name string) (inst interface{})
}

type InstanceFactory struct {
	instanceMap cmap.ConcurrentMap
}

func (f *InstanceFactory) Initialize(instanceMap cmap.ConcurrentMap)  {
	f.instanceMap = instanceMap
}

func (f *InstanceFactory) Initialized() bool  {
	return f.instanceMap != nil
}

func (f *InstanceFactory) SetInstance(name string, instance interface{}) {
	if _, ok := f.instanceMap.Get(name); ok && !gotest.IsRunning() {
		log.Fatalf("[factory] instance name %v is already taken", name)
	}
	f.instanceMap.Set(name, instance)
}

func (f *InstanceFactory) GetInstance(name string) (inst interface{}) {
	var ok bool
	if inst, ok = f.instanceMap.Get(name); !ok {
		return nil
	}
	return
}