package repos

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient struct {
	Username string
	Password string
	Host     string
}

type NameResult struct {
	Name      string
	Version   string
	Format    string
	ProjectID int `json:"project_id"`
}

func (h *HttpClient) NewRequest(method, endpoint string) (request *http.Request, err error) {
	url := fmt.Sprintf("%s/%s", h.Host, endpoint)
	request, err = http.NewRequest(method, url, nil)
	if err != nil {
		return
	}
	if h.Username != "" && h.Password != "" {
		request.SetBasicAuth(h.Username, h.Password)
	}
	return
}

func (h *HttpClient) Do(request *http.Request) (body []byte, resp *http.Response, err error) {

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err = client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err = ioutil.ReadAll(resp.Body)
		return
	}

	err = errors.New(resp.Status)
	return
}

func (h *HttpClient) http(method, endpoint string) ([]byte, *http.Response, error) {
	request, err := h.NewRequest(method, endpoint)
	if err != nil {
		return nil, nil, err
	}

	return h.Do(request)
}

func (h *HttpClient) Get(endpoint string) ([]byte, *http.Response, error) {
	return h.http(http.MethodGet, endpoint)
}

func (h *HttpClient) GetNameResult(url string) ([]NameResult, error) {
	body, _, err := h.Get(url)
	if err != nil {
		return nil, err
	}
	var result []NameResult
	if err = json.Unmarshal(body, &result); err != nil {
		return result, err
	}
	return result, err
}
