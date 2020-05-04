package events

import (
	"github.com/evgeny-klyopov/dep/app/events/models/config"
)

type setting struct {
	deployConfig *config.Config
	hostReleaseNumber *int
	variables map[string]string
}


func (e *event) setVariables() {
	host := e.deployConfig.Hosts[*(e.hostReleaseNumber)]
	e.variables = make(map[string]string)

	if e.deployConfig.Variables != nil {
		for k, v := range *(e.deployConfig.Variables) {
			e.variables[k] = v
		}
	}

	e.variables["{{release_path}}"] = host.DeployPath + "/release"
	e.variables["{{stage}}"] = e.stage
}
func (e *event) checkRepositoryConfig() error {
	host := e.deployConfig.Hosts[*(e.hostReleaseNumber)]

	if host.Branch == nil && e.deployConfig.Repository != nil {
		path := (*(e.config)).GetFilePathSetting()
		return e.getError(NoSetBranch, false, nil, e.stage, path)
	}

	return nil
}
func (e *event) setKeepRelease() {
	if e.deployConfig.KeepReleases == nil {
		defaultKeepRelease := (*(e.config)).GetDefaultKeepRelease()
		e.deployConfig.KeepReleases = &defaultKeepRelease
	}
}
func (e *event) setHost() error {
	for i, host := range e.deployConfig.Hosts {
		if host.Stage == e.stage {
			e.hostReleaseNumber = &i
			break
		}
	}

	if e.hostReleaseNumber == nil {
		path := (*(e.config)).GetFilePathSetting()
		return e.getError(NotFoundStage, false, nil, e.stage, path)
	}

	return nil
}
func (e *event) readDeployConfig() error {
	path := (*(e.config)).GetFilePathSetting()

	deployConfig, err := config.GetConfig(path)

	if err != nil {
		return err
	}

	e.setting.deployConfig = deployConfig

	return nil
}


type settinger interface {
	readDeployConfig() error
	checkRepositoryConfig() error
	setKeepRelease()
	setHost() error
	setVariables()
}