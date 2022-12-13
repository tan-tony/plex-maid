package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

func NewMusicbrainz() (*Musicbrainz, error) {
	u, _ := url.Parse("http://musicbrainz.org/ws/2/artist/")

	r := &Musicbrainz{
		Client:  http.Client{},
		baseURL: u,
		// https://musicbrainz.org/doc/MusicBrainz_API/Rate_Limiting
		limiter: rate.NewLimiter(rate.Every(1*time.Second), 1),
	}

	return r, nil
}

type Musicbrainz struct {
	baseURL *url.URL
	http.Client
	limiter *rate.Limiter
}

func (this *Musicbrainz) SearchArtist(ctx context.Context, name string) (*Artist, error) {
	this.limiter.Wait(ctx)

	u := *this.baseURL
	params := url.Values{
		"query":  []string{`"` + name + `"`},
		"limit":  []string{"-1"},
		"offset": []string{"-1"},
		"fmt":    []string{"json"},
		// switches the Solr query parser from edismax to dismax, which will escape certain special query syntax characters by default for ease of use.
		"dismax": []string{"false"},
	}
	u.RawQuery = params.Encode()

	resp, err := this.Get(u.String())
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 200:
	default:
		msg, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("Unexpected response code [%d]: %s", resp.StatusCode, msg)
	}

	type schema struct {
		Artists []struct {
			Name     string `json:"name"`
			SortName string `json:"sort-name"`
			Aliases  []struct {
				Name string `json:"name"`
			} `json:"aliases"`
		} `json;"artists"`
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	payload := schema{}
	if err := dec.Decode(&payload); err != nil {
		return nil, err
	}

	if len(payload.Artists) == 0 {
		return nil, ErrNotFound
	}

	pick := payload.Artists[0]
	artist := Artist{
		Name:     pick.Name,
		SortName: pick.SortName,
	}
	for _, alias := range pick.Aliases {
		artist.Aliases = append(artist.Aliases, alias.Name)
	}

	return &artist, nil
}
