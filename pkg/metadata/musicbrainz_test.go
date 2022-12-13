package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	musicbrainz, _ = NewMusicbrainz()
)

func TestMusicbrainzSearchArtist(t *testing.T) {
	cases := map[string]string{
		"jay chou":     "周杰倫",
		"jacky cheung": "張學友",
		"S.H.E":        "S.H.E",
		"F.I.R.":       "F.I.R.",
		"Winnie Hsin":  "辛曉琪",
	}

	for k, v := range cases {
		artist, err := musicbrainz.SearchArtist(k)
		if !assert.NoError(t, err) {
			t.Fatalf("Fatal when searching %s: %s", k, err)
		}
		assert.Equal(t, v, artist.Name)
	}
}
