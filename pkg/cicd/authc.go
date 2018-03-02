package cicd

type Authentication interface{
	Authenticate() error
}
