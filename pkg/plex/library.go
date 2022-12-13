package plex

import (
	"encoding/json"
)

type LibraryGetResponse struct {
	MediaContainer struct {
		Size      int `json:"size"`
		Directory []struct {
			// library ID
			Key        string `json:"key"`
			Title      string `json:"title"`
			Refreshing bool   `json:"refreshing"`
		} `json:"Directory"`
	} `json:"MediaContainer"`
}

func (this *Client) GetAllLibrary() (*LibraryGetResponse, error) {
	req := this.baseReq("GET", "/library/sections/all")

	resp, err := this.Do(req)
	if err := checkHTTPResponse(resp, err); err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data := LibraryGetResponse{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
