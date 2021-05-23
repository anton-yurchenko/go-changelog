package changelog_test

import (
	"testing"

	"github.com/anton-yurchenko/go-changelog"

	"github.com/stretchr/testify/assert"
)

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
