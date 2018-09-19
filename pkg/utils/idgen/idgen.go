package idgen

import (
	"fmt"
	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	sf = sonyflake.NewSonyflake(st)
	//if sf == nil {
	//	panic("sonyflake not created")
	//}
}

// Next generates next id as an uint64
func Next() (id uint64, err error) {
	var i uint64
	i, err = sf.NextID()
	if err == nil {
		id = i
	}
	return
}

// NextString generates next id as a string
func NextString() (id string, err error) {
	var i uint64
	i, err = sf.NextID()
	if err == nil {
		id = fmt.Sprintf("%d", i)
	}
	return
}
