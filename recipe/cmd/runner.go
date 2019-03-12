package main

import (
	"code.cloudfoundry.org/eirini"
	"code.cloudfoundry.org/eirini/recipe"
	"fmt"
	"os"
)

func main() {
	appID := os.Getenv(eirini.EnvAppID)
	stagingGUID := os.Getenv(eirini.EnvStagingGUID)
	completionCallback := os.Getenv(eirini.EnvCompletionCallback)
	eiriniAddress := os.Getenv(eirini.EnvEiriniAddress)
	appBitsDownloadURL := os.Getenv(eirini.EnvDownloadURL)
	dropletUploadURL := os.Getenv(eirini.EnvDropletUploadURL)
	buildpacksJSON := os.Getenv(eirini.EnvBuildpacks)

	downloadClient := createDownloadHTTPClient()

	installer := &recipe.PackageInstaller{
		Client:    downloadClient,
		Extractor: &recipe.Unzipper{},
	}

	uploader := &recipe.DropletUploader{
		HTTPClient: createUploaderHTTPClient(),
	}

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

	buildpacksKeyModifier := &recipe.BuildpacksKeyModifier{CCBuildpacksJSON: buildpacksJSON}

	executor := &recipe.PacksExecutor{
		Conf:           packsConf,
		Installer:      installer,
		Uploader:       uploader,
		Commander:      commander,
		ResultModifier: buildpacksKeyModifier,
	}

	recipeConf := recipe.Config{
		AppID:              appID,
		StagingGUID:        stagingGUID,
		CompletionCallback: completionCallback,
		EiriniAddr:         eiriniAddress,
		DropletUploadURL:   dropletUploadURL,
		PackageDownloadURL: appBitsDownloadURL,
	}

	err := executor.ExecuteRecipe(recipeConf)
	if err != nil {
		fmt.Println("Error while executing staging recipe:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Staging completed")
}
