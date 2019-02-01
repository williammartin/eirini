package recipe

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/eirini"
	"github.com/pkg/errors"
)

type PackageInstaller struct {
	Client    *http.Client
	Extractor eirini.Extractor
}

func (d *PackageInstaller) Install(downloadURL URL, targetDir string) error {
	if targetDir == "" {
		return errors.New("empty targetDir provided")
	}

	// To: Monday Mario & Steffen
	// consider using a TempDir to store the zip file
	// or if posible (preferably) unpack directly from the
	// download stream
	zipPath := filepath.Join(targetDir, appID) + ".zip"
	if err := d.download(downloadURL, zipPath); err != nil {
		return err
	}

	return d.Extractor.Extract(zipPath, targetDir)
}

func (d *PackageInstaller) download(appID string, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	appBits, err := d.fetchAppBits(appID)
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

func (d *PackageInstaller) fetchAppBits(appID string) (io.ReadCloser, error) {
	path, err := url.Parse("/v2/apps/" + appID + "/download")
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%s is not a legal app ID", appID))
	}

	url := d.ServerURL.ResolveReference(path)

	resp, err := d.Client.Get(url.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Download failed. Status Code %d", resp.StatusCode))
	}

	return resp.Body, nil
}
