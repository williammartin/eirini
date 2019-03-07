package recipe

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Buildpack struct {
	Name string
	Key  string
	Url  string
}

type BuildpackManager struct {
	Client *http.Client
}

func OpenBuildpackUrl(buildpack *Buildpack, client *http.Client) ([]byte, error) {
	resp, err := client.Get(buildpack.Url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request buildpack")
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
