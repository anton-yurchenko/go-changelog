package changelog_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/anton-yurchenko/go-changelog"
	"github.com/anton-yurchenko/go-changelog/mocks"
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
### Deprecated
- A

### Fixed
- A

## [0.0.1] - 2021-05-19
notice

### Added
- A
- B

### Changed
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

### Added
- A:
` + "```yaml" + `
this:
  that:
    - x` + "```" + `

- B

### Changed
- A
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
