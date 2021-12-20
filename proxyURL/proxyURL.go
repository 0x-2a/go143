package proxyURL

import (
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
	Do(req *http.Request) (resp *http.Response, err error)
}

type ProxyClient struct {
	httpClient HTTPClient
}

func NewProxyClient(httpClient HTTPClient) *ProxyClient {
	return &ProxyClient{
		httpClient: httpClient,
	}
}

func (r *ProxyClient) GetProxyURL(url string) []byte {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("Could not make req for googleBookRes, %s", err.Error())
		return nil
	}

	res, err := r.httpClient.Do(req)
	if err != nil {
		log.Errorf("Could not fetch %s, %s", url, err.Error())
		return nil
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Errorf("Non 200 status: %d, %s", res.StatusCode, url)
		return nil
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Errorf("Could not read polygon res: %s, %s", url, err.Error())
		return nil
	}

	return body
}
