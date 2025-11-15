package world

import "github.com/pelletier/go-toml"

const (
	InvalidOrEmpty = iota
)

type Invalid struct {
	Object World
	Error  error
}

func (i *Invalid) Invalid() error {
	_, err := toml.Marshal(&i.Object)
	return err
}
