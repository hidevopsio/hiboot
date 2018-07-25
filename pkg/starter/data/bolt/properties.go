package bolt

import "os"

type properties struct {
	Database string      `json:"database"`
	Mode     os.FileMode `json:"mode"`
	Timeout  int64       `json:"timeout"`
}
