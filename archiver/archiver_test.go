package archiver

import (
	"testing"
	"time"

	cache "github.com/patrickmn/go-cache"
)

type testMetadata struct {
	isRespost bool
	Response  Metadata
}

func TestFetch(t *testing.T) {
	cases := map[string]testMetadata{
		"https://imgs.xkcd.com/comics/shouldnt_be_hard.png": {
			false,
			Metadata{},
		},
		"https://imgs.xkcd.com/comics/bad_timing.png": {
			true,
			Metadata{
				User:    "garfield",
				Channel: "Techies",
			},
		},
	}

	testFixture := map[string]Metadata{
		"https://imgs.xkcd.com/comics/bad_timing.png": Metadata{
			User:    "garfield",
			Channel: "Techies",
		},
	}

	// Create cache with fixture
	createTestArchive(testFixture)
	defer archive.Flush()

	for k, tc := range cases {
		is_respost, data := archive.fetch(k)

		if tc.isRespost != is_respost {
			t.Errorf("Expected %+v got %+v", tc.isRespost, is_respost)
		}

		if tc.Response != data {
			t.Errorf("Expected %+v got %+v", tc.Response, data)
		}
	}
}

func TestCreate(t *testing.T) {
	testUrl := "https://imgs.xkcd.com/comics/shouldnt_be_hard.png"
	testUser := "garfield"
	testChannel := "Techies"

	expected := Metadata{
		User:    testUser,
		Channel: testChannel,
	}

	// Create empty cache
	Init(time.Duration(0) * time.Hour)
	defer archive.Flush()

	result := archive.create(testUrl, testUser, testChannel)
	if result != expected {
		t.Errorf("Wrong result back: Expected %+v but got %+v", result, expected)
	}

	isRepost, data := archive.fetch(testUrl)
	if isRepost != true {
		t.Error("Wrong key entry in cache: Expected this to be a repost")
	}
	if data != expected {
		t.Errorf("Wrong entry in cache: Expected %+v but got %+v", result, expected)
	}
}

func createTestArchive(items map[string]Metadata) {
	ti := make(map[string]cache.Item)
	for k, i := range items {
		ti[k] = cache.Item{
			Object:     i,
			Expiration: 0,
		}
	}

	archive.Cache = *cache.NewFrom(cache.NoExpiration, 10*time.Millisecond, ti)
}
