package executor

import (
	"github.com/hidevopsio/hi/boot/pkg/log"
	"fmt"
	"github.com/hidevopsio/hi/cicd/pkg/pipeline/pintf"
)

func Run(p pintf.PipelineInterface) error {

	p.EnsurePipeline()

	log.Debug(p)

	err := p.PullSourceCode()
	if err != nil {
		return fmt.Errorf("failed: %s", err)
	}

	err = p.Build()
	if err != nil {
		return fmt.Errorf("failed: %s", err)
	}

	err = p.RunUnitTest()
	if err != nil {
		return fmt.Errorf("failed: %s", err)
	}

	err = p.Deploy()
	if err != nil {
		return fmt.Errorf("failed: %s", err)
	}

	return nil
}
