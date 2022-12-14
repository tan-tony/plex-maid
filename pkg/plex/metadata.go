package plex

import (
	"encoding/json"
	"fmt"
)

type MetadataGetResponse struct {
	MediaContainer struct {
		// Size int `json:"size"`
		// TotalSize int `json:"totalSize"`
		Metadata []struct {
			RatingKey string `json:"ratingKey"`
			Title     string `json:"title"`
			UpdatedAt int64  `json:"updatedAt"`
			AddedAt   int64  `json:"addedAt"`
			Similar   []struct {
				ID  int    `json:"id"`
				Tag string `json:"tag"`
			} `json:"Similar"`
			Field []struct {
				Locked bool   `json:"locked"`
				Name   string `json:"name"`
			} `json:"Field"`
		} `json:"Metadata"`
	} `json:"MediaContainer"`
}

// GetMetadata from "library/metadata/{id}"
func (this *Client) GetMetadata(containerID string) (*MetadataGetResponse, error) {
	req := this.baseReq("GET", fmt.Sprintf("/library/metadata/%s", containerID))

	resp, err := this.Do(req)
	if err := checkHTTPResponse(resp, err); err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data := MetadataGetResponse{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
