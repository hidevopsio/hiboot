package pipeline

type Authentication interface{
	Authenticate() error
}
