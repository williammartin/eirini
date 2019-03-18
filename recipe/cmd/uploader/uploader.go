package main

import (
	"net/http"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/eirini"
	"code.cloudfoundry.org/eirini/recipe"
	"code.cloudfoundry.org/eirini/recipe/cmd/commons"
	"code.cloudfoundry.org/eirini/util"
)

func main() {
	dropletUploadURL := os.Getenv(eirini.EnvDropletUploadURL)

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

	responder := commons.NewResponder(cfg)

	client := createUploaderHTTPClient()
	err := uploadClient.Upload(client, stagerURL, dropletUploaderURL, commons.OutputDropletLocation)
	if err != nil {
		responder.RespondWithFailure(err)
		os.Exit(1)
	}

	buildpackCfg := os.Getenv(eirini.EnvBuildpacks)
	resp, err := responder.PrepareResponse(cfg, commons.OutputMetadataLocation, buildpackCfg)
	if err != nil {
		// TODO: log error
		commons.RespondWithFailure(cfg, err)
		os.Exit(1)
	}

	err = commons.RespondWithSuccess(resp)
	if err != nil {
		// TODO: log that it didnt go through
		os.Exit(1)
	}
}

func createUploaderHTTPClient() *http.Client {
	cert := filepath.Join(eirini.CCCertsMountPath, eirini.CCUploaderCertName)
	cacert := filepath.Join(eirini.CCCertsMountPath, eirini.CCInternalCACertName)
	key := filepath.Join(eirini.CCCertsMountPath, eirini.CCUploaderKeyName)

	client, err := util.CreateTLSHTTPClient([]util.CertPaths{
		{Crt: cert, Key: key, Ca: cacert},
	})
	if err != nil {
		panic(err)
	}

	return client
}
