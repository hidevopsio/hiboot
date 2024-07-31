package instantiate

import (
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/log"

	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
)

type instance struct {
	instMap cmap.ConcurrentMap
}

func newInstanceContainer(instMap cmap.ConcurrentMap) factory.InstanceContainer {
	if instMap == nil {
		instMap = cmap.New()
	}
	return &instance{
		instMap: instMap,
	}
}

// Get instanceContainer
func (i *instance) Get(params ...interface{}) (retVal interface{}) {
	name, obj := factory.ParseParams(params...)

	// get from instanceContainer map if external instanceContainer map does not have it
	if md, ok := i.instMap.Get(name); ok {
		metaData := factory.CastMetaData(md)
		if metaData != nil {
			switch obj.(type) {
			case factory.MetaData:
				retVal = metaData
			default:
				retVal = metaData.Instance
			}
		}
	}

	return
}

// Set save instanceContainer
func (i *instance) Set(params ...interface{}) (err error) {
	name, inst := factory.ParseParams(params...)

	metaData := factory.CastMetaData(inst)
	if metaData == nil {
		metaData = factory.NewMetaData(inst)
	}

	old, ok := i.instMap.Get(name)
	if ok {
		err = fmt.Errorf("instanceContainer %v is already taken by %v", name, old)
		// TODO: should handle such error
		log.Debugf("%+v", err)
		return
	}

	i.instMap.Set(name, metaData)
	return
}

// Items return map items
func (i *instance) Items() map[string]interface{} {
	return i.instMap.Items()
}
