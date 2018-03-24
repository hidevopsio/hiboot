
// dependencies: ci -> pipeline -> impl

package ci

import (
	"fmt"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"github.com/hidevopsio/hi/cicd/pkg/pipeline"
)

func Run(p pipeline.PipelineInterface) error {
	log.Debug("ci.Run()")
	err := p.EnsureParam()
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
