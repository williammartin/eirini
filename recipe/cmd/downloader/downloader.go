package main

import (
	"code.cloudfoundry.org/eirini"
	"code.cloudfoundry.org/eirini/recipe"
	"code.cloudfoundry.org/eirini/recipe/cmd/commons"
	"code.cloudfoundry.org/eirini/util"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

const buildPacksDir = "/var/lib/buildpacks"
const workspaceDir = "/workspace"

func main() {
	downloadClient := createDownloadHTTPClient()
	buildPackManager := NewManager(downloadClient, http.DefaultClient, buildPacksDir)

	installer := &recipe.PackageInstaller{
		Client:    downloadClient,
		Extractor: &recipe.Unzipper{},
	}

	var buildpacks []recipe.Buildpack
	err := json.Unmarshal([]byte(commons.BuildpackJson()), &buildpacks)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error unmarshaling environment variable %s: %s", eirini.EnvBuildpacks, err.Error()))
		os.Exit(1)
	}

	if err = buildPackManager.Install(buildpacks); err != nil {
		fmt.Println("Error while installing buildpacks:", err.Error())
		os.Exit(1)
	}

	recipeConf := commons.RecipeConfig()
	err = installer.Install(recipeConf.PackageDownloadURL, workspaceDir)
	if err != nil {
		fmt.Println("Error while installing app bits:", err.Error())
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
