package changelog_test

import (
	"testing"

	changelog "github.com/anton-yurchenko/go-changelog"

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

### Security
- A
- B

### Changed
- A
- B

### Added
- A
- B

### Removed
- A
- B

### Fixed
- A
- B

### Deprecated
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
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		a.Equal(test.Expected, test.Changelog.ToString())
	}
}

func TestAddNotice(t *testing.T) {
	a := assert.New(t)

	t.Log("Test Case 1/1 - Update Notice")

	target := &changelog.Changes{
		Notice: stringP("notice"),
	}

	expected := &changelog.Changes{
		Notice: stringP(""),
	}

	target.AddNotice("")

	a.Equal(expected, target)
}

func TestAddChange(t *testing.T) {
	a := assert.New(t)

	type expected struct {
		Changes *changelog.Changes
		Error   string
	}
	type test struct {
		Changes  *changelog.Changes
		Scope    string
		Change   string
		Expected expected
	}

	suite := map[string]test{
		"Empty": {
			Changes: new(changelog.Changes),
			Scope:   "Added",
			Change:  "",
			Expected: expected{
				Changes: new(changelog.Changes),
				Error:   "",
			},
		},
		"Append Added": {
			Changes: &changelog.Changes{
				Added: sliceOfStringsP([]string{"previous"}),
			},
			Scope:  "Added",
			Change: "change",
			Expected: expected{
				Changes: &changelog.Changes{
					Added: sliceOfStringsP([]string{"previous", "change"}),
				},
				Error: "",
			},
		},
		"Append Changed": {
			Changes: &changelog.Changes{
				Changed: sliceOfStringsP([]string{"previous"}),
			},
			Scope:  "Changed",
			Change: "change",
			Expected: expected{
				Changes: &changelog.Changes{
					Changed: sliceOfStringsP([]string{"previous", "change"}),
				},
				Error: "",
			},
		},
		"Append Deprecated": {
			Changes: &changelog.Changes{
				Deprecated: sliceOfStringsP([]string{"previous"}),
			},
			Scope:  "Deprecated",
			Change: "change",
			Expected: expected{
				Changes: &changelog.Changes{
					Deprecated: sliceOfStringsP([]string{"previous", "change"}),
				},
				Error: "",
			},
		},
		"Append Removed": {
			Changes: &changelog.Changes{
				Removed: sliceOfStringsP([]string{"previous"}),
			},
			Scope:  "Removed",
			Change: "change",
			Expected: expected{
				Changes: &changelog.Changes{
					Removed: sliceOfStringsP([]string{"previous", "change"}),
				},
				Error: "",
			},
		},
		"Append Fixed": {
			Changes: &changelog.Changes{
				Fixed: sliceOfStringsP([]string{"previous"}),
			},
			Scope:  "Fixed",
			Change: "change",
			Expected: expected{
				Changes: &changelog.Changes{
					Fixed: sliceOfStringsP([]string{"previous", "change"}),
				},
				Error: "",
			},
		},
		"Append Security": {
			Changes: &changelog.Changes{
				Security: sliceOfStringsP([]string{"previous"}),
			},
			Scope:  "Security",
			Change: "change",
			Expected: expected{
				Changes: &changelog.Changes{
					Security: sliceOfStringsP([]string{"previous", "change"}),
				},
				Error: "",
			},
		},
		"Added Scope": {
			Changes: new(changelog.Changes),
			Scope:   "Added",
			Change:  "change",
			Expected: expected{
				Changes: &changelog.Changes{
					Added: sliceOfStringsP([]string{"change"}),
				},
				Error: "",
			},
		},
		"Changed Scope": {
			Changes: new(changelog.Changes),
			Scope:   "Changed",
			Change:  "change",
			Expected: expected{
				Changes: &changelog.Changes{
					Changed: sliceOfStringsP([]string{"change"}),
				},
				Error: "",
			},
		},
		"Deprecated Scope": {
			Changes: new(changelog.Changes),
			Scope:   "Deprecated",
			Change:  "change",
			Expected: expected{
				Changes: &changelog.Changes{
					Deprecated: sliceOfStringsP([]string{"change"}),
				},
				Error: "",
			},
		},
		"Removed Scope": {
			Changes: new(changelog.Changes),
			Scope:   "Removed",
			Change:  "change",
			Expected: expected{
				Changes: &changelog.Changes{
					Removed: sliceOfStringsP([]string{"change"}),
				},
				Error: "",
			},
		},
		"Fixed Scope": {
			Changes: new(changelog.Changes),
			Scope:   "Fixed",
			Change:  "change",
			Expected: expected{
				Changes: &changelog.Changes{
					Fixed: sliceOfStringsP([]string{"change"}),
				},
				Error: "",
			},
		},
		"Security Scope": {
			Changes: new(changelog.Changes),
			Scope:   "Security",
			Change:  "change",
			Expected: expected{
				Changes: &changelog.Changes{
					Security: sliceOfStringsP([]string{"change"}),
				},
				Error: "",
			},
		},
		"Invalid Scope": {
			Changes: &changelog.Changes{
				Security: sliceOfStringsP([]string{"change"}),
			},
			Scope:  "Invalid",
			Change: "change",
			Expected: expected{
				Changes: &changelog.Changes{
					Security: sliceOfStringsP([]string{"change"}),
				},
				Error: "unexpected scope: Invalid (supported: [added,changed,deprecated,removed,fixed,security])",
			},
		},
		"Append Empty Added": {
			Changes: &changelog.Changes{
				Added: sliceOfStringsP([]string{"previous"}),
			},
			Scope:  "added",
			Change: "",
			Expected: expected{
				Changes: &changelog.Changes{
					Added: sliceOfStringsP([]string{"previous"}),
				},
				Error: "",
			},
		},
		"Append Empty Changed": {
			Changes: &changelog.Changes{
				Changed: sliceOfStringsP([]string{"previous"}),
			},
			Scope:  "changed",
			Change: "",
			Expected: expected{
				Changes: &changelog.Changes{
					Changed: sliceOfStringsP([]string{"previous"}),
				},
				Error: "",
			},
		},
		"Append Empty Deprecated": {
			Changes: &changelog.Changes{
				Deprecated: sliceOfStringsP([]string{"previous"}),
			},
			Scope:  "deprecated",
			Change: "",
			Expected: expected{
				Changes: &changelog.Changes{
					Deprecated: sliceOfStringsP([]string{"previous"}),
				},
				Error: "",
			},
		},
		"Append Empty Removed": {
			Changes: &changelog.Changes{
				Removed: sliceOfStringsP([]string{"previous"}),
			},
			Scope:  "removed",
			Change: "",
			Expected: expected{
				Changes: &changelog.Changes{
					Removed: sliceOfStringsP([]string{"previous"}),
				},
				Error: "",
			},
		},
		"Append Empty Fixed": {
			Changes: &changelog.Changes{
				Fixed: sliceOfStringsP([]string{"previous"}),
			},
			Scope:  "fixed",
			Change: "",
			Expected: expected{
				Changes: &changelog.Changes{
					Fixed: sliceOfStringsP([]string{"previous"}),
				},
				Error: "",
			},
		},
		"Append Empty Security": {
			Changes: &changelog.Changes{
				Security: sliceOfStringsP([]string{"previous"}),
			},
			Scope:  "security",
			Change: "",
			Expected: expected{
				Changes: &changelog.Changes{
					Security: sliceOfStringsP([]string{"previous"}),
				},
				Error: "",
			},
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		err := test.Changes.AddChange(test.Scope, test.Change)
		a.Equal(test.Expected.Changes, test.Changes)
		if test.Expected.Error != "" {
			a.EqualError(err, test.Expected.Error)
		}
	}
}
