package pipelines

type Pipeline interface {
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

