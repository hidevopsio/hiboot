package pintf

type PipelineInterface interface {
	EnsurePipeline() error
	PullSourceCode() error
	Build() error
	RunUnitTest() error
	RunIntegrationTest() error
	Analysis() error
	CopyTarget() error
	Upload() error
	NewImage() error
	Deploy() error
}

