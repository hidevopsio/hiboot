package bolt

import "os"

type Properties struct {
	Database string      `json:"database"`
	Mode     os.FileMode `json:"mode"`
	Timeout  int64       `json:"timeout"`
}
