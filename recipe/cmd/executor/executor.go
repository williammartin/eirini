package main

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/eirini"
	"code.cloudfoundry.org/eirini/recipe"
)

const (
	buildpacksDir             = "/var/lib/buildpacks"
	outputDropletLocation     = "/out/droplet.tgz"
	outputBuildArtifactsCache = "/cache/cache.tgz"
	outputMetadataLocation    = "/out/result.json"
)

func main() {

	stagingGUID := os.Getenv(eirini.EnvStagingGUID)
	completionCallback := os.Getenv(eirini.EnvCompletionCallback)
	eiriniAddress := os.Getenv(eirini.EnvEiriniAddress)

	responder := recipe.NewResponder(stagingGUID, completionCallback, eiriniAddress)

	commander := &recipe.IOCommander{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	packsConf := recipe.PacksBuilderConf{
		BuildpacksDir:             buildpacksDir,
		OutputDropletLocation:     outputDropletLocation,
		OutputBuildArtifactsCache: outputBuildArtifactsCache,
		OutputMetadataLocation:    outputMetadataLocation,
	}

	executor := &recipe.PacksExecutor{
		Conf:      packsConf,
		Commander: commander,
	}

	err := executor.ExecuteRecipe()
	if err != nil {
		responder.RespondWithFailure(err)
		os.Exit(1)
	}

	fmt.Println("Recipe Execution completed")
}
