package recipe

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Buildpack struct {
	Name string
	Key string
	Url string
	SkipDetect bool
}

type BuildpackInstaller struct {
	Client    *http.Client
}

func (b *BuildpackInstaller) OpenUrl(buildpack *Buildpack) ([]byte, error) {

	resp, err := b.Client.Get(buildpack.Url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("downloading buildpack failed with status code %d", resp.StatusCode))
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}