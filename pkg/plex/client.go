package plex

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

func NewClient(serverURL string) (*Client, error) {
	baseURL, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	r := &Client{
		baseURL: baseURL,
		Client:  http.Client{},
	}

	return r, nil
}

type Client struct {
	baseURL *url.URL
	http.Client
}

// returns the full url of a given request to plex
func (this *Client) baseReq(method string, urlPath string) *http.Request {
	baseURL := *this.baseURL
	baseURL.Path = path.Join(baseURL.Path, urlPath)
	req := &http.Request{
		Method: method,
		Header: http.Header{
			"Accept": []string{"application/json, text/plain, */*"},
		},
		URL: &baseURL,
	}

	return req
}

// checkHTTPResponse http client error and http response err
// reads response body in to error reason
func checkHTTPResponse(resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 200:
	default:
		msg, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return fmt.Errorf("Unexpected response code [%d]: %s", resp.StatusCode, msg)
	}

	return nil
}
