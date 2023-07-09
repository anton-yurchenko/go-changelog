package changelog_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/anton-yurchenko/go-changelog/mocks"

	changelog "github.com/anton-yurchenko/go-changelog"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func TestChangelogToString(t *testing.T) {
	a := assert.New(t)

	type test struct {
		Changelog *changelog.Changelog
		Expected  string
	}

	tm1, _ := time.Parse(changelog.DateFormat, "2021-05-19")
	tm2, _ := time.Parse(changelog.DateFormat, "2021-05-22")

	suite := map[string]test{
		"Empty": {
			Changelog: new(changelog.Changelog),
			Expected:  "",
		},
		"Title": {
			Changelog: &changelog.Changelog{
				Title: stringP("title"),
			},
			Expected: "# title\n",
		},
		"Description": {
			Changelog: &changelog.Changelog{
				Description: stringP("description\nhere"),
			},
			Expected: "description\nhere\n",
		},
		"Full Sorted": {
			Changelog: &changelog.Changelog{
				Title:       stringP("title"),
				Description: stringP("description\nhere"),
				Unreleased: &changelog.Release{
					URL: stringP("https://github.com/anton-yurchenko/go-changelog/compare/v0.0.2...HEAD"),
					Changes: &changelog.Changes{
						Added: sliceOfStringsP([]string{
							"A",
						}),
					},
				},
				Releases: changelog.Releases{
					{
						Version: stringP("0.0.1"),
						URL:     stringP("https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1"),
						Date:    &tm1,
						Changes: &changelog.Changes{
							Notice: stringP("notice"),
							Added: sliceOfStringsP([]string{
								"A",
								"B",
							}),
							Changed: sliceOfStringsP([]string{
								"A",
								"B",
							}),
						},
					},
					{
						Version: stringP("0.0.2"),
						URL:     stringP("https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...v0.0.2"),
						Date:    &tm2,
						Changes: &changelog.Changes{
							Fixed: sliceOfStringsP([]string{
								"A",
							}),
							Deprecated: sliceOfStringsP([]string{
								"A",
							}),
						},
					},
				},
			},
			Expected: `# title

description
here

## [Unreleased]

### Added

- A

## [0.0.2] - 2021-05-22

### Fixed

- A

### Deprecated

- A

## [0.0.1] - 2021-05-19

notice

### Changed

- A
- B

### Added

- A
- B

[Unreleased]: https://github.com/anton-yurchenko/go-changelog/compare/v0.0.2...HEAD
[0.0.2]: https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1`,
		},
		"Complex": {
			Changelog: &changelog.Changelog{
				Title:       stringP("title"),
				Description: stringP("description\nhere"),
				Releases: changelog.Releases{
					{
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
				},
			},
			Expected: `# title

description
here

## [0.0.1] - 2021-05-19

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

[0.0.1]: https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1`,
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		a.Equal(test.Expected, test.Changelog.ToString())
	}
}

func TestSaveToFile(t *testing.T) {
	a := assert.New(t)

	type test struct {
		Changelog *changelog.Changelog
		Expected  string
	}

	tm1, _ := time.Parse(changelog.DateFormat, "2021-05-19")
	tm2, _ := time.Parse(changelog.DateFormat, "2021-05-22")

	suite := map[string]test{
		"Empty": {
			Changelog: new(changelog.Changelog),
			Expected:  "reason",
		},
		"Full Sorted": {
			Changelog: &changelog.Changelog{
				Title:       stringP("title"),
				Description: stringP("description\nhere"),
				Unreleased: &changelog.Release{
					URL: stringP("https://github.com/anton-yurchenko/go-changelog/compare/v0.0.2...HEAD"),
					Changes: &changelog.Changes{
						Added: sliceOfStringsP([]string{
							"A",
						}),
					},
				},
				Releases: changelog.Releases{
					{
						Version: stringP("0.0.1"),
						URL:     stringP("https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1"),
						Date:    &tm1,
						Changes: &changelog.Changes{
							Notice: stringP("notice"),
							Added: sliceOfStringsP([]string{
								"A",
								"B",
							}),
							Changed: sliceOfStringsP([]string{
								"A",
								"B",
							}),
						},
					},
					{
						Version: stringP("0.0.2"),
						URL:     stringP("https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...v0.0.2"),
						Date:    &tm2,
						Changes: &changelog.Changes{
							Fixed: sliceOfStringsP([]string{
								"A",
							}),
							Deprecated: sliceOfStringsP([]string{
								"A",
							}),
						},
					},
				},
			},
			Expected: "reason",
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		m := new(mocks.Filesystem)
		if test.Expected != "" {
			m.On("Create", "CHANGELOG.md").Return(nil, errors.New(test.Expected)).Once()
		}

		err := test.Changelog.SaveToFile(m, "CHANGELOG.md")
		if test.Expected != "" {
			a.EqualError(err, fmt.Sprintf("error creating a file: %v", test.Expected))
		} else {
			a.Equal(nil, err)
		}
	}
}

func TestNewChangelog(t *testing.T) {
	a := assert.New(t)

	t.Log("Test Case 1/1 - Create New Changelog")

	expected := &changelog.Changelog{
		Releases: []*changelog.Release{},
	}

	a.Equal(expected, changelog.NewChangelog())
}

func TestSetTitle(t *testing.T) {
	a := assert.New(t)

	t.Log("Test Case 1/1 - Update Title")

	target := &changelog.Changelog{
		Title: stringP("changes"),
	}

	expected := &changelog.Changelog{
		Title: stringP("changelog"),
	}

	target.SetTitle("changelog")

	a.Equal(expected, target)
}

func TestSetDescription(t *testing.T) {
	a := assert.New(t)

	t.Log("Test Case 1/1 - Update Description")

	target := &changelog.Changelog{
		Description: stringP("description"),
	}

	expected := &changelog.Changelog{
		Description: stringP(""),
	}

	target.SetDescription("")

	a.Equal(expected, target)
}

func TestSetUnreleasedURL(t *testing.T) {
	a := assert.New(t)

	type test struct {
		Changelog *changelog.Changelog
		URL       string
		Expected  string
	}

	suite := map[string]test{
		"Valid URL": {
			Changelog: new(changelog.Changelog),
			URL:       "https://github.com/anton-yurchenko/go-changelog",
			Expected:  "",
		},
		"Invalid URL": {
			Changelog: new(changelog.Changelog),
			URL:       "github.com/sdf\as",
			Expected:  "parse \"github.com/sdf\\as\": net/url: invalid control character in URL",
		},
		"Replace URL": {
			Changelog: &changelog.Changelog{
				Unreleased: &changelog.Release{
					URL: stringP("gitlab.com"),
				},
			},
			URL:      "github.com",
			Expected: "",
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		if test.Expected != "" {
			a.EqualError(test.Changelog.SetUnreleasedURL(test.URL), test.Expected)
		} else {
			a.Equal(nil, test.Changelog.SetUnreleasedURL(test.URL))
		}
	}
}

func TestAddUnreleasedChange(t *testing.T) {
	a := assert.New(t)

	type expected struct {
		Changelog *changelog.Changelog
		Error     string
	}
	type test struct {
		Changelog *changelog.Changelog
		Scope     string
		Change    string
		Expected  expected
	}

	suite := map[string]test{
		"Valid": {
			Changelog: new(changelog.Changelog),
			Scope:     "Fixed",
			Change:    "change",
			Expected: expected{
				Changelog: &changelog.Changelog{
					Unreleased: &changelog.Release{
						Changes: &changelog.Changes{
							Fixed: sliceOfStringsP([]string{"change"}),
						},
					},
				},
				Error: "",
			},
		},
		"Invalid": {
			Changelog: new(changelog.Changelog),
			Scope:     "Invalid",
			Change:    "change",
			Expected: expected{
				Changelog: &changelog.Changelog{
					Unreleased: new(changelog.Release),
				},
				Error: "unexpected scope: Invalid (supported: [added,changed,deprecated,removed,fixed,security])",
			},
		},
		"Missing Unreleased Field": {
			Changelog: new(changelog.Changelog),
			Scope:     "Fixed",
			Change:    "",
			Expected: expected{
				Changelog: &changelog.Changelog{
					Unreleased: &changelog.Release{
						Changes: new(changelog.Changes),
					},
				},
				Error: "",
			},
		},
		"Missing Changes Field": {
			Changelog: &changelog.Changelog{
				Unreleased: new(changelog.Release),
			},
			Scope:  "Fixed",
			Change: "",
			Expected: expected{
				Changelog: &changelog.Changelog{
					Unreleased: &changelog.Release{
						Changes: new(changelog.Changes),
					},
				},
				Error: "",
			},
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		err := test.Changelog.AddUnreleasedChange(test.Scope, test.Change)

		if test.Expected.Error != "" {
			a.EqualError(err, test.Expected.Error)
		} else {
			a.Equal(test.Expected.Changelog, test.Changelog)
		}
	}
}

func TestChangelogGetRelease(t *testing.T) {
	a := assert.New(t)

	type test struct {
		Changelog *changelog.Changelog
		Release   *changelog.Release
		Version   string
	}

	suite := map[string]test{
		"Found": {
			Changelog: &changelog.Changelog{
				Releases: changelog.Releases{
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

		r := test.Changelog.GetRelease(test.Version)
		a.Equal(test.Release, r)
	}
}

func TestChangelogCreateRelease(t *testing.T) {
	a := assert.New(t)

	type expected struct {
		Release *changelog.Release
		Error   string
	}

	type test struct {
		Changelog *changelog.Changelog
		Version   string
		Date      string
		Expected  expected
	}

	suite := map[string]test{
		"Success": {
			Changelog: &changelog.Changelog{
				Releases: changelog.Releases{
					{
						Version: stringP("0.0.1"),
					},
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
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		r, err := test.Changelog.CreateRelease(test.Version, test.Date)
		if test.Expected.Error != "" {
			a.EqualError(err, test.Expected.Error)
		} else {
			a.Equal(test.Expected.Release, r)
		}
	}
}

func TestChangelogCreateReleaseWithURL(t *testing.T) {
	a := assert.New(t)

	type expected struct {
		Release *changelog.Release
		Error   string
	}

	type test struct {
		Changelog *changelog.Changelog
		Version   string
		Date      string
		URL       string
		Expected  expected
	}

	suite := map[string]test{
		"Success": {
			Changelog: &changelog.Changelog{
				Releases: changelog.Releases{
					{
						Version: stringP("0.0.1"),
					},
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
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		r, err := test.Changelog.CreateReleaseWithURL(test.Version, test.Date, test.URL)
		if test.Expected.Error != "" {
			a.EqualError(err, test.Expected.Error)
		} else {
			a.Equal(test.Expected.Release, r)
		}
	}
}

func TestCreateReleaseFromUnreleased(t *testing.T) {
	a := assert.New(t)

	type expected struct {
		Release *changelog.Release
		Error   string
	}

	type test struct {
		Changelog *changelog.Changelog
		Version   string
		Date      string
		Expected  expected
	}

	suite := map[string]test{
		"Success": {
			Changelog: &changelog.Changelog{
				Releases: changelog.Releases{
					{
						Version: stringP("0.0.1"),
					},
				},
				Unreleased: &changelog.Release{
					Changes: &changelog.Changes{
						Added:   sliceOfStringsP([]string{"feature"}),
						Changed: sliceOfStringsP([]string{"behavior"}),
					},
				},
			},
			Version: "1.0.0",
			Date:    "2021-05-30",
			Expected: expected{
				Release: &changelog.Release{
					Version: stringP("1.0.0"),
					Date:    parseDate("2021-05-30"),
					Changes: &changelog.Changes{
						Added:   sliceOfStringsP([]string{"feature"}),
						Changed: sliceOfStringsP([]string{"behavior"}),
					},
				},
				Error: "",
			},
		},
		"Failure": {
			Changelog: &changelog.Changelog{
				Releases: changelog.Releases{
					{
						Version: stringP("1.0.0"),
					},
				},
				Unreleased: &changelog.Release{
					Changes: &changelog.Changes{
						Added:   sliceOfStringsP([]string{"feature"}),
						Changed: sliceOfStringsP([]string{"behavior"}),
					},
				},
			},
			Version: "1.0.0",
			Date:    "2021-05-30",
			Expected: expected{
				Release: nil,
				Error:   "version 1.0.0 already exists",
			},
		},
		"Missing Changes": {
			Changelog: &changelog.Changelog{
				Releases: changelog.Releases{
					{
						Version: stringP("0.0.1"),
					},
				},
				Unreleased: &changelog.Release{},
			},
			Version: "1.0.0",
			Date:    "2021-05-30",
			Expected: expected{
				Release: nil,
				Error:   "missing 'Unreleased' section",
			},
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		r, err := test.Changelog.CreateReleaseFromUnreleased(test.Version, test.Date)
		if test.Expected.Error != "" {
			a.EqualError(err, test.Expected.Error)
		} else {
			a.Equal(test.Expected.Release, r)
			a.Equal(new(changelog.Release), test.Changelog.Unreleased)
		}
	}
}

func TestCreateReleaseFromUnreleasedWithURL(t *testing.T) {
	a := assert.New(t)

	type expected struct {
		Release *changelog.Release
		Error   string
	}

	type test struct {
		Changelog *changelog.Changelog
		Version   string
		Date      string
		URL       string
		Expected  expected
	}

	suite := map[string]test{
		"Success": {
			Changelog: &changelog.Changelog{
				Releases: changelog.Releases{
					{
						Version: stringP("0.0.1"),
					},
				},
				Unreleased: &changelog.Release{
					Changes: &changelog.Changes{
						Added:   sliceOfStringsP([]string{"feature"}),
						Changed: sliceOfStringsP([]string{"behavior"}),
					},
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
					Changes: &changelog.Changes{
						Added:   sliceOfStringsP([]string{"feature"}),
						Changed: sliceOfStringsP([]string{"behavior"}),
					},
				},
				Error: "",
			},
		},
		"Failure": {
			Changelog: &changelog.Changelog{
				Releases: changelog.Releases{
					{
						Version: stringP("1.0.0"),
					},
				},
				Unreleased: &changelog.Release{
					Changes: &changelog.Changes{
						Added:   sliceOfStringsP([]string{"feature"}),
						Changed: sliceOfStringsP([]string{"behavior"}),
					},
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
			Changelog: &changelog.Changelog{
				Releases: changelog.Releases{
					{
						Version: stringP("0.1.0"),
					},
				},
				Unreleased: &changelog.Release{
					Changes: &changelog.Changes{
						Added:   sliceOfStringsP([]string{"feature"}),
						Changed: sliceOfStringsP([]string{"behavior"}),
					},
				},
			},
			Version: "1.0.0",
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

		r, err := test.Changelog.CreateReleaseFromUnreleasedWithURL(test.Version, test.Date, test.URL)
		if test.Expected.Error != "" {
			a.EqualError(err, test.Expected.Error)
		} else {
			a.Equal(test.Expected.Release, r)
			a.Equal(new(changelog.Release), test.Changelog.Unreleased)
		}
	}
}
