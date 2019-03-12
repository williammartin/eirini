package main

import (
	"code.cloudfoundry.org/eirini"
	"code.cloudfoundry.org/eirini/util"
	"net/http"
	"path/filepath"
)

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
