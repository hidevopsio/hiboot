package pipelines

type Authentication interface{
	Authenticate() error
}
