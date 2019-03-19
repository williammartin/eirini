package main

import (
	"net/http"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/eirini"
	"code.cloudfoundry.org/eirini/recipe"
	"code.cloudfoundry.org/eirini/util"
)

const (
	OutputMetadataLocation = "/out/result.json"
	OutputDropletLocation  = "/out/droplet.tgz"
)

func main() {
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
	client := createUploaderHTTPClient()
	uploadClient := recipe.DropletUploader{
		Client: client,
	}

	err := uploadClient.Upload(dropletUploadURL, OutputDropletLocation)
	if err != nil {
		responder.RespondWithFailure(err)
		os.Exit(1)
	}

	buildpackCfg := os.Getenv(eirini.EnvBuildpacks)
	resp, err := responder.PrepareSuccessResponse(OutputMetadataLocation, buildpackCfg)
	if err != nil {
		// TODO: log error
		responder.RespondWithFailure(err)
		os.Exit(1)
	}

	err = responder.RespondWithSuccess(resp)
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
