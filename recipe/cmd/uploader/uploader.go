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
	buildpackCfg := os.Getenv(eirini.EnvBuildpacks)
	stagingGUID := os.Getenv(eirini.EnvStagingGUID)
	completionCallback := os.Getenv(eirini.EnvCompletionCallback)
	eiriniAddress := os.Getenv(eirini.EnvEiriniAddress)
	dropletUploadURL := os.Getenv(eirini.EnvDropletUploadURL)

	responder := recipe.NewResponder(stagingGUID, completionCallback, eiriniAddress)

	client, err := createUploaderHTTPClient()
	if err != nil {
		responder.RespondWithFailure(err)
		os.Exit(1)
	}

	uploadClient := recipe.DropletUploader{
		Client: client,
	}

	err = uploadClient.Upload(dropletUploadURL, OutputDropletLocation)
	if err != nil {
		responder.RespondWithFailure(err)
		os.Exit(1)
	}

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

func createUploaderHTTPClient() (*http.Client, error) {
	cert := filepath.Join(eirini.CCCertsMountPath, eirini.CCUploaderCertName)
	cacert := filepath.Join(eirini.CCCertsMountPath, eirini.CCInternalCACertName)
	key := filepath.Join(eirini.CCCertsMountPath, eirini.CCUploaderKeyName)

	return util.CreateTLSHTTPClient([]util.CertPaths{
		{Crt: cert, Key: key, Ca: cacert},
	})
}
