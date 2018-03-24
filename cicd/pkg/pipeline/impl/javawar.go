package impl

import "github.com/hidevopsio/hi/boot/pkg/log"

type JavaWarPipeline struct{
	JavaPipeline
}

func (p *JavaWarPipeline) Deploy() error {
	log.Debug("JavaWarPipeline.Deploy()")
	return nil
}