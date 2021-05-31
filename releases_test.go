package changelog_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/anton-yurchenko/go-changelog"

	"github.com/stretchr/testify/assert"
)

func parseDate(date string) *time.Time {
	t, err := time.Parse(changelog.DateFormat, date)
	if err != nil {
		return nil
	}

	return &t
}

func TestReleasesSwap(t *testing.T) {
	a := assert.New(t)

	t.Log("Test Case 1/1 - Releases Swap")

	releases := changelog.Releases{
		{
			Version: stringP("1.1.1"),
		},
		{
			Version: stringP("1.2.2"),
		},
	}

	expected := changelog.Releases{
		{
			Version: stringP("1.2.2"),
		},
		{
			Version: stringP("1.1.1"),
		},
	}

	releases.Swap(0, 1)
	a.Equal(expected, releases)
}

func TestReleasesLen(t *testing.T) {
	a := assert.New(t)

	t.Log("Test Case 1/1 - Releases Len")

	releases := changelog.Releases{
		{
			Version: stringP("1.1.1"),
		},
		{
			Version: stringP("1.2.2"),
		},
		{
			Version: stringP("1.0.0"),
		},
	}

	a.Equal(3, releases.Len())
}

func TestReleasesLess(t *testing.T) {
	a := assert.New(t)

	type test struct {
		Releases changelog.Releases
		Expected bool
	}

	suite := map[string]test{
		"Smaller": {
			Releases: changelog.Releases{
				{
					Version: stringP("1.0.0"),
				},
				{
					Version: stringP("2.0.0"),
				},
			},
			Expected: true,
		},
		"Equal": {
			Releases: changelog.Releases{
				{
					Version: stringP("1.0.0"),
				},
				{
					Version: stringP("1.0.0"),
				},
			},
			Expected: false,
		},
		"Bigger": {
			Releases: changelog.Releases{
				{
					Version: stringP("2.0.0"),
				},
				{
					Version: stringP("1.0.0"),
				},
			},
			Expected: false,
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		a.Equal(test.Expected, test.Releases.Less(0, 1))
	}
}

func TestReleasesGetRelease(t *testing.T) {
	a := assert.New(t)

	type test struct {
		Releases *changelog.Releases
		Release  *changelog.Release
		Version  string
	}

	suite := map[string]test{
		"Not Found": {
			Releases: new(changelog.Releases),
			Version:  "1.0.0",
		},
		"Found": {
			Releases: &changelog.Releases{
				{
					Version: stringP("1.1.0"),
				},
				{
					Version: stringP("1.0.0"),
				},
				{
					Version: stringP("0.1.0"),
				},
			},
			Release: &changelog.Release{
				Version: stringP("1.0.0"),
			},
			Version: "1.0.0",
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		r := test.Releases.GetRelease(test.Version)
		a.Equal(test.Release, r)
	}
}

func TestReleasesCreateRelease(t *testing.T) {
	a := assert.New(t)

	type expected struct {
		Release *changelog.Release
		Error   string
	}

	type test struct {
		Releases *changelog.Releases
		Version  string
		Date     string
		Expected expected
	}

	suite := map[string]test{
		"Success": {
			Releases: &changelog.Releases{
				{
					Version: stringP("0.0.1"),
				},
			},
			Version: "1.0.0",
			Date:    "2021-05-30",
			Expected: expected{
				Release: &changelog.Release{
					Version: stringP("1.0.0"),
					Date:    parseDate("2021-05-30"),
					Changes: new(changelog.Changes),
				},
				Error: "",
			},
		},
		"Existing Version": {
			Releases: &changelog.Releases{
				{
					Version: stringP("1.0.0"),
				},
			},
			Version: "1.0.0",
			Date:    "2021-05-30",
			Expected: expected{
				Release: nil,
				Error:   "version 1.0.0 already exists",
			},
		},
		"Invalid Timestamp": {
			Releases: &changelog.Releases{
				{
					Version: stringP("0.1.0"),
				},
			},
			Version: "1.0.0",
			Date:    "202020-0505-3030",
			Expected: expected{
				Release: nil,
				Error:   fmt.Sprintf("invalid date 202020-0505-3030, expected to match regex %v", changelog.DateRegex),
			},
		},
		"Invalid Version": {
			Releases: &changelog.Releases{
				{
					Version: stringP("1.0.0"),
				},
			},
			Version: "1.0",
			Date:    "2021-05-30",
			Expected: expected{
				Release: nil,
				Error:   fmt.Sprintf("invalid semantic version 1.0, expected to match regex %v", changelog.SemVerRegex),
			},
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		r, err := test.Releases.CreateRelease(test.Version, test.Date)
		if test.Expected.Error != "" {
			a.EqualError(err, test.Expected.Error)
		} else {
			a.Equal(test.Expected.Release, r)
		}
	}
}

func TestReleasesCreateReleaseWithURL(t *testing.T) {
	a := assert.New(t)

	type expected struct {
		Release *changelog.Release
		Error   string
	}

	type test struct {
		Releases *changelog.Releases
		Version  string
		Date     string
		URL      string
		Expected expected
	}

	suite := map[string]test{
		"Success": {
			Releases: &changelog.Releases{
				{
					Version: stringP("0.0.1"),
				},
			},
			Version: "1.0.0",
			Date:    "2021-05-30",
			URL:     "https://github.com/anton-yurchenko/go-changelog",
			Expected: expected{
				Release: &changelog.Release{
					Version: stringP("1.0.0"),
					Date:    parseDate("2021-05-30"),
					URL:     stringP("https://github.com/anton-yurchenko/go-changelog"),
					Changes: new(changelog.Changes),
				},
				Error: "",
			},
		},
		"Existing Version": {
			Releases: &changelog.Releases{
				{
					Version: stringP("1.0.0"),
				},
			},
			Version: "1.0.0",
			Date:    "2021-05-30",
			URL:     "https://github.com/anton-yurchenko/go-changelog",
			Expected: expected{
				Release: nil,
				Error:   "version 1.0.0 already exists",
			},
		},
		"Invalid URL": {
			Releases: &changelog.Releases{
				{
					Version: stringP("1.0.0"),
				},
			},
			Version: "2.0.0",
			Date:    "2021-05-30",
			URL:     "github.com/sdf\as",
			Expected: expected{
				Release: nil,
				Error:   "parse \"github.com/sdf\\as\": net/url: invalid control character in URL",
			},
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		r, err := test.Releases.CreateReleaseWithURL(test.Version, test.Date, test.URL)
		if test.Expected.Error != "" {
			a.EqualError(err, test.Expected.Error)
		} else {
			a.Equal(test.Expected.Release, r)
		}
	}
}
