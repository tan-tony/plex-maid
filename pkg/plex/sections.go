package plex

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/go-querystring/query"
)

type LibraryType = int

var (
	ByArtist LibraryType = 8
	ByAlbum  LibraryType = 9
	ByTrack  LibraryType = 10
)

type SectionsGetRequest struct {
	Type           LibraryType `url:"type"`
	ContainerStart int         `url:"X-Plex-Container-Start"`
	ContainerSize  int         `url:"X-Plex-Container-Size"`
}

type SectionsGetResponse struct {
	MediaContainer struct {
		Size      int `json:"size"`
		TotalSize int `json:"totalSize"`
		Metadata  []struct {
			RatingKey string `json:"ratingKey"`
			Title     string `json:"title"`
		} `json:"Metadata"`
	} `json:"MediaContainer"`
}

type SectionsPutRequest struct {
	Type            LibraryType `url:"type"`
	ID              string      `url:"id"`
	Title           string      `url:"title.value"`
	TitleSort       string      `url:"titleSort.value"`
	TitleLocked     int         `url:"title.locked"`
	TitleSortLocked int         `url:"titleSort.locked"`
}

// GetAllSections for a given library from "library/sections/{id}/all"
func (this *Client) GetAllSections(ctx context.Context, libraryID string, libraryType LibraryType, c chan<- *SectionsGetResponse) error {
	go func() {
		var (
			batchSize = 5
			head      = 0
		)

		defer close(c)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				req := this.baseReq("GET", fmt.Sprintf("/library/sections/%s/all", libraryID))
				params := SectionsGetRequest{
					Type:           libraryType,
					ContainerStart: head,
					ContainerSize:  batchSize,
				}
				q, _ := query.Values(params)
				req.URL.RawQuery = q.Encode()

				resp, err := this.Do(req)
				if err := checkHTTPResponse(resp, err); err != nil {
					log.Printf("Failed to fetch library sections: %s", err)
					continue
				}

				data := SectionsGetResponse{}
				dec := json.NewDecoder(resp.Body)
				if err := dec.Decode(&data); err != nil {
					log.Printf("Fatal error on trying to read library sections: %s", err)
					resp.Body.Close()
					continue
				}
				resp.Body.Close()

				// send to channel
				select {
				case <-ctx.Done():
				case c <- &data:
				}

				if data.MediaContainer.Size < batchSize {
					// done
					return
				} else {
					head += batchSize
				}
			}
		}
	}()

	return nil
}

// PutSections to "/library/sections/{id}/all"
func (this *Client) PutSections(libraryID string, params SectionsPutRequest) error {
	req := this.baseReq("PUT", fmt.Sprintf("/library/sections/%s/all", libraryID))
	q, _ := query.Values(params)
	req.URL.RawQuery = q.Encode()

	resp, err := this.Do(req)
	return checkHTTPResponse(resp, err)
}
