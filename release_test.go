package changelog_test

import (
	"fmt"
	"testing"
	"time"

	changelog "github.com/anton-yurchenko/go-changelog"

	"github.com/stretchr/testify/assert"
)

func TestReleaseToString(t *testing.T) {
	a := assert.New(t)

	type test struct {
		Release            *changelog.Release
		ExpectedString     string
		ExpectedDefinition string
	}

	tm1, _ := time.Parse(changelog.DateFormat, "2021-05-19")

	suite := map[string]test{
		"Unreleased": {
			Release:            new(changelog.Release),
			ExpectedString:     "## [Unreleased]\n",
			ExpectedDefinition: "",
		},
		"Unreleased With URL": {
			Release: &changelog.Release{
				URL: stringP("https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...HEAD"),
			},
			ExpectedString:     "## [Unreleased]\n",
			ExpectedDefinition: "[Unreleased]: https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...HEAD",
		},
		"Release": {
			Release: &changelog.Release{
				Version: stringP("0.0.1"),
				URL:     stringP("https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1"),
				Date:    &tm1,
			},
			ExpectedString:     "## [0.0.1] - 2021-05-19\n",
			ExpectedDefinition: "[0.0.1]: https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1",
		},
		"Without Date": {
			Release: &changelog.Release{
				Version: stringP("0.0.1"),
				URL:     stringP("https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1"),
			},
			ExpectedString:     "## [0.0.1]\n",
			ExpectedDefinition: "[0.0.1]: https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1",
		},
		"Complex Release with Changes": {
			Release: &changelog.Release{
				Version: stringP("0.0.1"),
				URL:     stringP("https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1"),
				Date:    &tm1,
				Changes: &changelog.Changes{
					Notice: stringP("notice"),
					Added: sliceOfStringsP([]string{
						"A:\n```yaml\nthis:\n  that:\n    - x```\n",
						"B",
					}),
					Changed: sliceOfStringsP([]string{
						"A",
						"B",
					}),
				},
			},
			ExpectedString: `## [0.0.1] - 2021-05-19

notice

### Changed

- A
- B

### Added

- A:
` + "```yaml" + `
this:
  that:
    - x` + "```" + `

- B
`,
			ExpectedDefinition: "[0.0.1]: https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1",
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		r, l := test.Release.ToString()
		a.Equal(test.ExpectedString, r)
		a.Equal(test.ExpectedDefinition, l)
	}
}

func TestSetDate(t *testing.T) {
	a := assert.New(t)

	type expected struct {
		Release *changelog.Release
		Error   string
	}

	type test struct {
		Release  *changelog.Release
		Date     string
		Expected expected
	}

	suite := map[string]test{
		"Valid": {
			Release: new(changelog.Release),
			Date:    "2020-05-30",
			Expected: expected{
				Release: &changelog.Release{
					Date: parseDate("2020-05-30"),
				},
				Error: "",
			},
		},
		"Invalid": {
			Release: new(changelog.Release),
			Date:    "202020-0505-3030",
			Expected: expected{
				Release: nil,
				Error:   fmt.Sprintf("invalid date 202020-0505-3030, expected to match regex %v", changelog.DateRegex),
			},
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		err := test.Release.SetDate(test.Date)
		if test.Expected.Error != "" {
			a.EqualError(err, test.Expected.Error)
		} else {
			a.Equal(nil, err)
		}
	}
}

func TestSetVersion(t *testing.T) {
	a := assert.New(t)

	type expected struct {
		Release *changelog.Release
		Error   string
	}

	type test struct {
		Release  *changelog.Release
		Version  string
		Expected expected
	}

	suite := map[string]test{
		"Valid": {
			Release: new(changelog.Release),
			Version: "1.0.0",
			Expected: expected{
				Release: &changelog.Release{
					Version: stringP("1.0.0"),
				},
				Error: "",
			},
		},
		"Invalid": {
			Release: new(changelog.Release),
			Version: "1.0",
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

		err := test.Release.SetVersion(test.Version)
		if test.Expected.Error != "" {
			a.EqualError(err, test.Expected.Error)
		} else {
			a.Equal(nil, err)
		}
	}
}

func TestSetURL(t *testing.T) {
	a := assert.New(t)

	type expected struct {
		Release *changelog.Release
		Error   string
	}

	type test struct {
		Release  *changelog.Release
		URL      string
		Expected expected
	}

	suite := map[string]test{
		"Valid": {
			Release: new(changelog.Release),
			URL:     "https://github.com/anton-yurchenko/go-changelog",
			Expected: expected{
				Release: &changelog.Release{
					URL: stringP("https://github.com/anton-yurchenko/go-changelog"),
				},
				Error: "",
			},
		},
		"Invalid": {
			Release: new(changelog.Release),
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

		err := test.Release.SetURL(test.URL)
		if test.Expected.Error != "" {
			a.EqualError(err, test.Expected.Error)
		} else {
			a.Equal(nil, err)
		}
	}
}

func TestReleaseAddNotice(t *testing.T) {
	a := assert.New(t)

	t.Log("Test Case 1/1 - Update Notice")

	target := &changelog.Release{
		Changes: &changelog.Changes{
			Notice: stringP("notice"),
		},
	}

	expected := &changelog.Release{
		Changes: &changelog.Changes{
			Notice: stringP(""),
		},
	}

	target.AddNotice("")

	a.Equal(expected, target)
}

func TestReleaseAddChange(t *testing.T) {
	a := assert.New(t)

	t.Log("Test Case 1/1 - Add Change")

	target := &changelog.Release{
		Changes: &changelog.Changes{},
	}

	expected := &changelog.Release{
		Changes: &changelog.Changes{
			Added: sliceOfStringsP([]string{"change"}),
		},
	}

	err := target.AddChange("added", "change")

	a.Equal(expected, target)
	a.Equal(nil, err)
}
