package types

type Handler interface {
	Name() string
	Help() string
	Open() error
	Do() error
	Close() error
}

