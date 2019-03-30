package recipe

import bap "code.cloudfoundry.org/buildpackapplifecycle"

type Executor interface {
	ExecuteRecipe() error
}

//go:generate counterfeiter . StagingResultModifier
type StagingResultModifier interface {
	Modify(result bap.StagingResult) (bap.StagingResult, error)
}

//go:generate counterfeiter . Uploader
type Uploader interface {
	Upload(path, url string) error
}

//go:generate counterfeiter . Installer
type Installer interface {
	Install() error
}

//go:generate counterfeiter . Commander
type Commander interface {
	Exec(cmd string, args ...string) error
}
