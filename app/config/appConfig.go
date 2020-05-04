package config

import (
	"github.com/evgeny-klyopov/bashColor"
)

type metaSetting struct {
	name      string
	separator string
}

type appConfig struct {
	abstract
	filePathSetting   string
	defaultKeepRelease   int8
	deployArchiveName string
	metaSetting       *metaSetting
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

type AppConfigurator interface {
	Configurator
	GetFilePathSetting() string
	GetDeployArchiveName() string
	GetMetaSetting() *metaSetting
	GetDefaultKeepRelease() int8
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
			name:      "meta",
			separator: ", ",
		},
	}
}
