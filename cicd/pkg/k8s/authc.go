package k8s

type Authentication interface{
	Authenticate() error
}
