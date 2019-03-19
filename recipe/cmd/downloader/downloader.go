package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/eirini"
	"code.cloudfoundry.org/eirini/recipe"
	"code.cloudfoundry.org/eirini/util"
)

const buildPacksDir = "/var/lib/buildpacks"
const workspaceDir = "/workspace"

func main() {
	downloadClient := createDownloadHTTPClient()
	buildPackManager := recipe.NewBuildpackManager(downloadClient, http.DefaultClient, buildPacksDir)

	installer := &recipe.PackageInstaller{
		Client:    downloadClient,
		Extractor: &recipe.Unzipper{},
	}

	appID := os.Getenv(eirini.EnvAppID)
	stagingGUID := os.Getenv(eirini.EnvStagingGUID)
	completionCallback := os.Getenv(eirini.EnvCompletionCallback)
	eiriniAddress := os.Getenv(eirini.EnvEiriniAddress)
	appBitsDownloadURL := os.Getenv(eirini.EnvDownloadURL)
	dropletUploadURL := os.Getenv(eirini.EnvDropletUploadURL)
	buildpacksJSON := os.Getenv(eirini.EnvBuildpacks)

	cfg := recipe.Config{
		AppID:              appID,
		StagingGUID:        stagingGUID,
		CompletionCallback: completionCallback,
		EiriniAddr:         eiriniAddress,
		DropletUploadURL:   dropletUploadURL,
		PackageDownloadURL: appBitsDownloadURL,
	}

	responder := recipe.NewResponder(cfg)

	var buildpacks []recipe.Buildpack
	err := json.Unmarshal([]byte(buildpacksJSON), &buildpacks)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error unmarshaling environment variable %s: %s", eirini.EnvBuildpacks, err.Error()))
		responder.RespondWithFailure(err)
		os.Exit(1)
	}

	if err = buildPackManager.Install(buildpacks); err != nil {
		fmt.Println("Error while installing buildpacks:", err.Error())
		responder.RespondWithFailure(err)
		os.Exit(1)
	}

	err = installer.Install(cfg.PackageDownloadURL, workspaceDir)
	if err != nil {
		fmt.Println("Error while installing app bits:", err.Error())
		responder.RespondWithFailure(err)
		os.Exit(1)
	}

	fmt.Println("Downloading completed")
}

func createDownloadHTTPClient() *http.Client {
	apiCA := filepath.Join(eirini.CCCertsMountPath, eirini.CCInternalCACertName)
	cert := filepath.Join(eirini.CCCertsMountPath, eirini.CCAPICertName)
	key := filepath.Join(eirini.CCCertsMountPath, eirini.CCAPIKeyName)

	client, err := util.CreateTLSHTTPClient([]util.CertPaths{
		{Crt: cert, Key: key, Ca: apiCA},
	})

	if err != nil {
		panic(err)
	}

	return client
}
