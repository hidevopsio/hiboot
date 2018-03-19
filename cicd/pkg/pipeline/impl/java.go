package pipelines

import (
	"github.com/hidevopsio/hi/boot/pkg/log"
	"github.com/hidevopsio/hi/cicd/pkg/pipeline"
)

type JavaPipeline struct{
	pipeline.Pipeline
}

func (p *JavaPipeline) Deploy() error {
	log.Debug("JavaPipeline.Deploy()")
	return nil
}