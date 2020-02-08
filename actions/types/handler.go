package types

type Handler interface {
	Name() string
	Help() string
	GetEnvironment() error
	Check() error
	Connect() error
	Do() error
}

