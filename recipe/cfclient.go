package recipe

import (
	"net/http"

	"github.com/pkg/errors"
)

type CfClient struct {
	httpClient *http.Client
}

func (c *CfClient) GetDropletByAppGuid(string) ([]byte, error) {
	return nil, nil
}

func (c *CfClient) PushDroplet(string, string) error {
	return errors.New("Implement me!")
}

func (c *CfClient) GetAppBitsByAppGuid(string) (*http.Response, error) {
	return nil, nil
}
