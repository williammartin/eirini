package recipe

import (
	"os"
	"os/exec"
)

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
	PacksBuilderPath          string
	BuildpacksDir             string
	OutputDropletLocation     string
	OutputBuildArtifactsCache string
	OutputMetadataLocation    string
}

type PacksExecutor struct {
	Conf      PacksBuilderConf
	Commander Commander
}

func (e *PacksExecutor) ExecuteRecipe() error {
	args := []string{
		"-buildpacksDir", e.Conf.BuildpacksDir,
		"-outputDroplet", e.Conf.OutputDropletLocation,
		"-outputBuildArtifactsCache", e.Conf.OutputBuildArtifactsCache,
		"-outputMetadata", e.Conf.OutputMetadataLocation,
	}

	err := e.Commander.Exec(e.Conf.PacksBuilderPath, args...)
	if err != nil {
		return err
	}
	return nil
}
