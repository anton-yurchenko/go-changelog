package changelog_test

import (
	"testing"

	"github.com/anton-yurchenko/go-changelog"

	"github.com/stretchr/testify/assert"
)

func TestChangesToString(t *testing.T) {
	a := assert.New(t)

	type test struct {
		Changelog *changelog.Changes
		Expected  string
	}

	suite := map[string]test{
		"Empty": {
			Changelog: new(changelog.Changes),
			Expected:  "",
		},
		"Notice": {
			Changelog: &changelog.Changes{
				Notice: stringP("notice"),
			},
			Expected: "notice\n",
		},
		"Added": {
			Changelog: &changelog.Changes{
				Added: sliceOfStringsP([]string{
					"A",
					"B",
				}),
			},
			Expected: "### Added\n- A\n- B\n",
		},
		"Changed": {
			Changelog: &changelog.Changes{
				Changed: sliceOfStringsP([]string{
					"A",
					"B",
				}),
			},
			Expected: "### Changed\n- A\n- B\n",
		},
		"Deprecated": {
			Changelog: &changelog.Changes{
				Deprecated: sliceOfStringsP([]string{
					"A",
					"B",
				}),
			},
			Expected: "### Deprecated\n- A\n- B\n",
		},
		"Removed": {
			Changelog: &changelog.Changes{
				Removed: sliceOfStringsP([]string{
					"A",
					"B",
				}),
			},
			Expected: "### Removed\n- A\n- B\n",
		},
		"Fixed": {
			Changelog: &changelog.Changes{
				Fixed: sliceOfStringsP([]string{
					"A",
					"B",
				}),
			},
			Expected: "### Fixed\n- A\n- B\n",
		},
		"Security": {
			Changelog: &changelog.Changes{
				Security: sliceOfStringsP([]string{
					"A",
					"B",
				}),
			},
			Expected: "### Security\n- A\n- B\n",
		},
		"Full": {
			Changelog: &changelog.Changes{
				Notice: stringP("notice"),
				Added: sliceOfStringsP([]string{
					"A",
					"B",
				}),
				Changed: sliceOfStringsP([]string{
					"A",
					"B",
				}),
				Deprecated: sliceOfStringsP([]string{
					"A",
					"B",
				}),
				Removed: sliceOfStringsP([]string{
					"A",
					"B",
				}),
				Fixed: sliceOfStringsP([]string{
					"A",
					"B",
				}),
				Security: sliceOfStringsP([]string{
					"A",
					"B",
				}),
			},
			Expected: `notice

### Added
- A
- B

### Changed
- A
- B

### Deprecated
- A
- B

### Removed
- A
- B

### Fixed
- A
- B

### Security
- A
- B
`,
		},
		"Multiline": {
			Changelog: &changelog.Changes{
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
			Expected: `notice

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
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		a.Equal(test.Expected, test.Changelog.ToString())
	}
}
