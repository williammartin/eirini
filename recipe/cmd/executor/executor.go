package main

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/eirini"
	"code.cloudfoundry.org/eirini/recipe"
)

func main() {
	commander := &recipe.IOCommander{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	packsConf := recipe.PacksBuilderConf{
		BuildpacksDir:             "/var/lib/buildpacks",
		OutputDropletLocation:     "/out/droplet.tgz",
		OutputBuildArtifactsCache: "/cache/cache.tgz",
		OutputMetadataLocation:    "/out/result.json",
	}
	executor := &recipe.PacksExecutor{
		Conf:      packsConf,
		Commander: commander,
	}

	appID := os.Getenv(eirini.EnvAppID)
	stagingGUID := os.Getenv(eirini.EnvStagingGUID)
	completionCallback := os.Getenv(eirini.EnvCompletionCallback)
	eiriniAddress := os.Getenv(eirini.EnvEiriniAddress)
	appBitsDownloadURL := os.Getenv(eirini.EnvDownloadURL)
	dropletUploadURL := os.Getenv(eirini.EnvDropletUploadURL)

	cfg := recipe.Config{
		AppID:              appID,
		StagingGUID:        stagingGUID,
		CompletionCallback: completionCallback,
		EiriniAddr:         eiriniAddress,
		DropletUploadURL:   dropletUploadURL,
		PackageDownloadURL: appBitsDownloadURL,
	}

	responder := recipe.NewResponder(cfg)

	err := executor.ExecuteRecipe(cfg)
	if err != nil {
		responder.RespondWithFailure(err)
		os.Exit(1)
	}

	fmt.Println("Recipe Execution completed")
}
