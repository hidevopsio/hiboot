package inject

import "reflect"

type valueTag struct {
	BaseTag
}


func init() {
	AddTag("value", new(valueTag))
}


func (t *valueTag) IsSingleton() bool  {
	return false
}

func (t *valueTag) Decode(object reflect.Value, field reflect.StructField, tag string) (retVal interface{}) {
	if tag != "" {
		//log.Debug(valueTag)
		retVal = t.replaceReferences(tag)
	}
	return retVal
}
