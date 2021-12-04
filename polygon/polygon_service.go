package polygon

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	polygonAPIBase = "https://api.polygon.io/"
)

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
	Do(req *http.Request) (resp *http.Response, err error)
}

type RestClient struct {
	polygonAPIKey string
	httpClient    HTTPClient
}

func NewRestClient(polygonAPIKey string, httpClient HTTPClient) *RestClient {
	return &RestClient{
		polygonAPIKey: polygonAPIKey,
		httpClient:    httpClient,
	}
}

func (r *RestClient) GetPolygonPath(path string) []byte {
	polygonURL := fmt.Sprintf("%s%s", polygonAPIBase, r.withToken(path))

	req, err := http.NewRequest(http.MethodGet, polygonURL, nil)
	if err != nil {
		log.Errorf("Could not make req for googleBookRes, %s", err.Error())
		return nil
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := r.httpClient.Do(req)
	if err != nil {
		log.Errorf("Could not fetch %s, %s", path, err.Error())
		return nil
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Errorf("Non 200 status: %d, %s", res.StatusCode, path)
		return nil
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Errorf("Could not read polygon res: %s, %s", path, err.Error())
		return nil
	}

	return body
}

func (r *RestClient) withToken(url string) string {
	urlWithToken := fmt.Sprintf("%s?apiKey=%s", url, r.polygonAPIKey)
	if strings.Contains(url, "?") {
		urlWithToken = fmt.Sprintf("%s&apiKey=%s", url, r.polygonAPIKey)
	}

	return urlWithToken
}
