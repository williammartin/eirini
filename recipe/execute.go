package recipe

import (
	"os"
	"os/exec"
)

const workspaceDir = "/workspace"

type IOCommander struct {
	Stdout *os.File
	Stderr *os.File
	Stdin  *os.File
}

func (c *IOCommander) Exec(cmd string, args ...string) error {
	command := exec.Command(cmd, args...) //#nosec
	command.Stdout = c.Stdout
	command.Stderr = c.Stderr
	command.Stdin = c.Stdin

	return command.Run()
}

type PacksBuilderConf struct {
	BuildpacksDir             string
	OutputDropletLocation     string
	OutputBuildArtifactsCache string
	OutputMetadataLocation    string
}

type Config struct {
	AppID              string
	StagingGUID        string
	CompletionCallback string
	EiriniAddr         string
	DropletUploadURL   string
	PackageDownloadURL string
}

type PacksExecutor struct {
	Conf      PacksBuilderConf
	Commander Commander
}

func (e *PacksExecutor) ExecuteRecipe(recipeConf Config) error {
	args := []string{
		"-buildpacksDir", e.Conf.BuildpacksDir,
		"-outputDroplet", e.Conf.OutputDropletLocation,
		"-outputBuildArtifactsCache", e.Conf.OutputBuildArtifactsCache,
		"-outputMetadata", e.Conf.OutputMetadataLocation,
	}

	err := e.Commander.Exec("/packs/builder", args...)
	if err != nil {
		return err
	}
	return nil
}
