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

	uploader := &recipe.DropletUploader{
		HTTPClient: createUploaderHTTPClient(),
	}

	err := uploader.Upload(commons.OutputDropletLocation, dropletUploadURL)
	if err != nil {
		commons.RespondWithFailure(err)
		os.Exit(1)
	}

	err = commons.RespondWithSuccess()
	if err != nil {
		commons.RespondWithFailure(err)
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
