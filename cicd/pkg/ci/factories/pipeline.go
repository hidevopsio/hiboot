package factories

import (
	"github.com/hidevopsio/hi/boot/pkg/log"
	"errors"
	"fmt"
	"github.com/hidevopsio/hi/cicd/pkg/ci"
	"github.com/hidevopsio/hi/cicd/pkg/ci/impl"
)

type PipelineFactory struct{}

const (
	JavaPipelineType = "java"
	JavaWarPipelineType = "java-war"
	NodeJsPipelineType = "nodejs"
	GitbookPipelineType = "gitbook"
)

func (pf *PipelineFactory) New(pipelineType string) (ci.PipelineInterface, error) {
	log.Debug("pf.NewPipeline()")
	switch pipelineType {
	case JavaPipelineType:
		return new(impl.JavaPipeline), nil
	case JavaWarPipelineType:
		return new(impl.JavaWarPipeline), nil
	default:
		return nil, errors.New(fmt.Sprintf("pipeline type %d not recognized\n", pipelineType))
	}
	return nil, nil
}
