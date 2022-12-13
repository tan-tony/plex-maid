package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"

	"plex-maid/pkg/metadata"
	"plex-maid/pkg/plex"
)

var (
	pc        *plex.Client
	scrobbler metadata.Scrobbler
)

func main() {
	flagPlexURL := flag.String("url", "", "URL of plex")
	flagToLock := flag.Bool("toLock", false, "Lock the change in Plex")
	flag.Parse()

	var (
		err error
	)

	pc, err = plex.NewClient(*flagPlexURL)
	if err != nil {
		log.Fatalf("Failed to start Plex client: %s", err)
	}

	scrobbler, err = metadata.NewMusicbrainz()
	if err != nil {
		log.Fatalf("Failed to start scrobbler: %s", err)
	}

	resp, err := pc.GetAllLibrary()
	if err != nil {
		log.Fatalf("Failed to fetch all library: %s", err)
	}
	log.Printf("Found library: %+v", resp.MediaContainer.Directory)

	for _, lib := range resp.MediaContainer.Directory {
		if !lib.Refreshing {
			scanAndUpdateLibrary(lib.Key, *flagToLock)
		}
	}
}

func scanAndUpdateLibrary(libraryID string, lock bool) {
	sectionsC := make(chan *plex.SectionsGetResponse)
	ctx := context.Background()

	err := pc.GetAllSections(ctx, libraryID, plex.ByArtist, sectionsC)
	if err != nil {
		log.Fatal(err)
	}

	for containers := range sectionsC {
		for _, obj := range containers.MediaContainer.Metadata {
			name := obj.Title
			log.Printf("Found %s on plex: %+v", name, obj)

			if locked, err := isArtistLocked(obj.RatingKey); err != nil {
				log.Printf("Failed to determine if artist is locked: %s", err)
			} else if locked {
				log.Printf("%s is locked, skipping", name)
				continue
			}

			artist, err := scrobbler.SearchArtist(ctx, name)
			if err != nil {
				log.Printf("Error while trying to search artist: %s", err)
				continue
			}
			artist.ToSimplified()
			log.Printf("Musicbrainz found match for %s: %+v", name, artist)

			boolToInt := func(b bool) int {
				if b {
					return 1
				}
				return 0
			}

			// update artist with new metadata
			// enhance searching experience by adding aliases in title sort
			titleSort := fmt.Sprintf("%s | %s", artist.SortName, artist.Aliases)
			params := plex.SectionsPutRequest{
				Type:            plex.ByArtist,
				ID:              obj.RatingKey,
				Title:           artist.Name,
				TitleSort:       titleSort,
				TitleLocked:     boolToInt(lock),
				TitleSortLocked: boolToInt(lock),
			}
			if err := pc.PutSections(libraryID, params); err != nil {
				log.Printf("Failed to update artist %s: %s", name, err)
			}
		}
	}
}

func isArtistLocked(id string) (bool, error) {
	// check if lock
	artistContainer, err := pc.GetMetadata(id)
	if err != nil {
		log.Printf("Failed to fetch metadata of %s: %s", id, err)
		return false, err
	}

	if len(artistContainer.MediaContainer.Metadata) == 0 {
		return false, errors.New("Unexpected metadata length 0")
	}

	for _, field := range artistContainer.MediaContainer.Metadata[0].Field {
		if field.Name == "title" && field.Locked {
			return true, nil
		}
	}

	return false, nil
}
