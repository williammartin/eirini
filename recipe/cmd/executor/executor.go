package main

import (
	"code.cloudfoundry.org/eirini/recipe"
	"code.cloudfoundry.org/eirini/recipe/cmd/commons"
	"fmt"
	"os"
)

func main() {
	commander := &recipe.IOCommander{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	packsConf := commons.PacksConfig()
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

	err := executor.ExecuteRecipe(cfg)
	if err != nil {
		commons.RespondWithFailure(err)
		os.Exit(1)
	}

	fmt.Println("Recipe Execution completed")
}
