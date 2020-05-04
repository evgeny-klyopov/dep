package config

import (
	"github.com/evgeny-klyopov/bashColor"
)

type abstract struct {
	version             string
	color             *bashColor.Colorer
}

type Configurator interface {
	GetVersion() string
	GetColor() *bashColor.Colorer
}

func (a *abstract) GetVersion() string {
	return a.version
}
func (a *abstract) GetColor() *bashColor.Colorer {
	return a.color
}

