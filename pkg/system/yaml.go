package system

import (
	"bytes"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	"strings"
)

func ReadYamlFromFile(file string) (prop map[string]interface{}, err error) {
	fs := afero.NewOsFs()
	var fb []byte
	fb, err = afero.ReadFile(fs, file)
	if err == nil {
		buf := new(bytes.Buffer)
		r := bytes.NewReader(fb)
		buf.ReadFrom(r)
		str := string(buf.Bytes())
		s := strings.Split(str, "---")
		var src []byte
		if len(s) == 1 {
			src = buf.Bytes()
		} else {
			src = []byte(s[1])
		}

		err = yaml.Unmarshal(src, &prop)
	}
	return
}
