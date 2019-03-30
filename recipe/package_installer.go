package recipe

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"code.cloudfoundry.org/eirini"
	"github.com/pkg/errors"
)

type PackageInstaller struct {
	client      *http.Client
	extractor   eirini.Extractor
	downloadURL string
	extractDir  string
}

func NewPackageInstaller(client *http.Client, extractor eirini.Extractor, downloadURL, extractDir string) Installer {
	return &PackageInstaller{
		client:      client,
		extractor:   extractor,
		downloadURL: downloadURL,
		extractDir:  extractDir,
	}
}

func (d *PackageInstaller) Install() error {
	if d.downloadURL == "" {
		return errors.New("empty downloadURL provided")
	}

	if d.extractDir == "" {
		return errors.New("empty extractDir provided")
	}

	zipPath := "/tmp/app.zip"

	err := d.download(d.downloadURL, zipPath)
	if err != nil {
		return err
	}

	return d.extractor.Extract(zipPath, d.extractDir)
}

func (d *PackageInstaller) download(downloadURL string, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	appBits, err := d.fetchAppBits(downloadURL)
	if err != nil {
		return err
	}

	defer appBits.Close()

	_, err = io.Copy(file, appBits)
	if err != nil {
		return errors.Wrap(err, "failed to copy content to file")
	}

	return nil
}

func (d *PackageInstaller) fetchAppBits(url string) (io.ReadCloser, error) {
	resp, err := d.client.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Download failed. Status Code %d", resp.StatusCode))
	}

	return resp.Body, nil
}
