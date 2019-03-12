package main

import (
	"code.cloudfoundry.org/eirini"
	"code.cloudfoundry.org/eirini/recipe"
	"code.cloudfoundry.org/eirini/recipe/cmd/buildpack"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const buildPacksDir = "/var/lib/buildpacks"
const workspaceDir = "/workspace"

func main() {
	appID := os.Getenv(eirini.EnvAppID)
	stagingGUID := os.Getenv(eirini.EnvStagingGUID)
	completionCallback := os.Getenv(eirini.EnvCompletionCallback)
	eiriniAddress := os.Getenv(eirini.EnvEiriniAddress)
	appBitsDownloadURL := os.Getenv(eirini.EnvDownloadURL)
	dropletUploadURL := os.Getenv(eirini.EnvDropletUploadURL)
	buildpacksJSON := os.Getenv(eirini.EnvBuildpacks)

	downloadClient := createDownloadHTTPClient()
	buildPackManager := buildpack.New(downloadClient, http.DefaultClient, buildPacksDir)

	installer := &recipe.PackageInstaller{
		Client:    downloadClient,
		Extractor: &recipe.Unzipper{},
	}

	recipeConf := recipe.Config{
		AppID:              appID,
		StagingGUID:        stagingGUID,
		CompletionCallback: completionCallback,
		EiriniAddr:         eiriniAddress,
		DropletUploadURL:   dropletUploadURL,
		PackageDownloadURL: appBitsDownloadURL,
	}

	var buildpacks []recipe.Buildpack
	err := json.Unmarshal([]byte(buildpacksJSON), &buildpacks)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error unmarshaling environment variable %s: %s", eirini.EnvBuildpacks, err.Error()))
		os.Exit(1)
	}

	if err = buildPackManager.Install(buildpacks); err != nil {
		fmt.Println("Error while installing buildpacks:", err.Error())
		os.Exit(1)
	}

	err = installer.Install(recipeConf.PackageDownloadURL, workspaceDir)
	if err != nil {
		fmt.Println("Error while installing app bits:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Downloading completed")
}
