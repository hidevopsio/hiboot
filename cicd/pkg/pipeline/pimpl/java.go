package pipelines

import "github.com/hidevopsio/hi/boot/pkg/log"

type JavaPipeline struct{
	Pipeline
}

func (p *JavaPipeline) Deploy() error {
	log.Debug("JavaPipeline.Deploy()")
	return nil
}