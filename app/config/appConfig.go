package config

import (
	"github.com/evgeny-klyopov/bashColor"
)

type metaSetting struct {
	Name      string
	Separator string
}

type appConfig struct {
	abstract
	filePathSetting   string
	defaultKeepRelease   int8
	deployArchiveName string
	metaSetting       *metaSetting
	flags flags
}

func (a *appConfig) GetFilePathSetting() string {
	return a.filePathSetting
}

func (a *appConfig) GetDeployArchiveName() string {
	return a.deployArchiveName
}

func (a *appConfig) GetMetaSetting() *metaSetting {
	return a.metaSetting
}

func (a *appConfig) GetDefaultKeepRelease() int8 {
	return a.defaultKeepRelease
}
func (a *appConfig) SetDebugFlag(value bool)  {
	a.flags.Debug = value
}
func (a *appConfig) GetFlags() flags {
	return a.flags
}

type AppConfigurator interface {
	Configurator
	SetDebugFlag(value bool)
	GetFilePathSetting() string
	GetDeployArchiveName() string
	GetMetaSetting() *metaSetting
	GetDefaultKeepRelease() int8
	GetFlags() flags
}

type flags struct{
	Debug bool
}

func NewAppConfig(version string, color *bashColor.Colorer) AppConfigurator {
	return &appConfig{
		abstract: abstract{
			version: version,
			color:   color,
		},
		filePathSetting:   "deploy.json",
		defaultKeepRelease:   10,
		deployArchiveName: "deploy.tar.gz",
		metaSetting: &metaSetting{
			Name:      "meta",
			Separator: ", ",
		},
	}
}
