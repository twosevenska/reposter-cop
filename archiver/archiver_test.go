package archiver

import (
	"testing"
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

	for k, tc := range cases {
		is_respost, data := fetch(k)

		if tc.isRespost != is_respost {
			t.Errorf("Expected %+v got %+v", tc.isRespost, is_respost)
		}

		if tc.Response != data {
			t.Errorf("Expected %+v got %+v", tc.Response, data)
		}
	}
}
