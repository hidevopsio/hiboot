package data

type Configuration interface {
	NewRepository(name string)
}
