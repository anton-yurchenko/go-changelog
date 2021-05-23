package changelog_test

import (
	"testing"
	"time"

	"github.com/anton-yurchenko/go-changelog"

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
			ExpectedString:     "## [Unreleased]",
			ExpectedDefinition: "",
		},
		"Unreleased With URL": {
			Release: &changelog.Release{
				URL: stringP("https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...HEAD"),
			},
			ExpectedString:     "## [Unreleased]",
			ExpectedDefinition: "[Unreleased]: https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...HEAD",
		},
		"Release": {
			Release: &changelog.Release{
				Version: stringP("0.0.1"),
				URL:     stringP("https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1"),
				Date:    &tm1,
			},
			ExpectedString:     "## [0.0.1] - 2021-05-19",
			ExpectedDefinition: "[0.0.1]: https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1",
		},
		"Without Date": {
			Release: &changelog.Release{
				Version: stringP("0.0.1"),
				URL:     stringP("https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1"),
			},
			ExpectedString:     "## [0.0.1]",
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
