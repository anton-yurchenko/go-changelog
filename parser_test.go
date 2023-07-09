package changelog_test

import (
	"testing"
	"time"

	changelog "github.com/anton-yurchenko/go-changelog"
	"github.com/anton-yurchenko/go-changelog/mocks"

	"github.com/pkg/errors"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func stringP(v string) *string {
	return &v
}

func sliceOfStringsP(v []string) *[]string {
	return &v
}

func dateP(v string) *time.Time {
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return nil
	}

	return &t
}

func TestNewParserWithFilesystem(t *testing.T) {
	a := assert.New(t)
	roBase := afero.NewReadOnlyFs(afero.NewOsFs())
	fs := afero.NewCopyOnWriteFs(roBase, afero.NewMemMapFs())

	type expected struct {
		Result *changelog.Parser
		Error  string
	}

	type test struct {
		Filesystem afero.Fs
		File       string
		Expected   expected
	}

	suite := map[string]test{
		"Correct": {
			Filesystem: fs,
			File:       "CHANGELOG1.md",
			Expected: expected{
				Result: &changelog.Parser{
					Filepath:   "CHANGELOG1.md",
					Filesystem: fs,
				},
				Error: "",
			},
		},
		"Missing File": {
			Filesystem: fs,
			File:       "CHANGELOG2.md",
			Expected: expected{
				Result: nil,
				Error:  "file CHANGELOG2.md not found: stat CHANGELOG2.md: no such file or directory",
			},
		},
	}

	var counter int
	for name, test := range suite {
		counter++
		t.Logf("Test Case %v/%v - %s", counter, len(suite), name)

		// prepare test case
		switch name {
		case "Missing File":
		default:
			if err := afero.WriteFile(test.Filesystem, test.File, []byte("# Changelog"), 0644); err != nil {
				t.Errorf("error preparing test case: error creating file %v: %v", test.File, err)
				continue
			}
		}

		// test
		p, err := changelog.NewParserWithFilesystem(test.Filesystem, test.File)
		a.Equal(test.Expected.Result, p)
		if test.Expected.Error != "" || err != nil {
			a.EqualError(err, test.Expected.Error)
		}

		// cleanup
		switch name {
		case "Missing File":
		default:
			if err := test.Filesystem.Remove(test.File); err != nil {
				t.Errorf("error cleanup: error removing file %v: %v", test.File, err)
			}
		}
	}
}

func TestNewParser(t *testing.T) {
	a := assert.New(t)

	expectedResult := &changelog.Parser{
		Filepath:   "CHANGELOG.md",
		Filesystem: afero.NewOsFs(),
	}

	p, err := changelog.NewParser("CHANGELOG.md")
	a.Equal(expectedResult, p)
	a.Equal(nil, err)
}

func TestParser(t *testing.T) {
	a := assert.New(t)
	roBase := afero.NewReadOnlyFs(afero.NewOsFs())
	fs := afero.NewCopyOnWriteFs(roBase, afero.NewMemMapFs())

	type expected struct {
		Result *changelog.Changelog
		Error  string
	}

	type test struct {
		Changelog string
		Expected  expected
	}

	suite := map[string]test{
		"Empty": {
			Changelog: "",
			Expected: expected{
				Result: &changelog.Changelog{
					Releases: []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Filesystem Error": {
			Changelog: "",
			Expected: expected{
				Error: "error loading a buffer: reason",
			},
		},
		"Title": {
			Changelog: "# Changelog",
			Expected: expected{
				Result: &changelog.Changelog{
					Title:    stringP("Changelog"),
					Releases: []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Title - Empty Trailing Line": {
			Changelog: `# Changelog
`,
			Expected: expected{
				Result: &changelog.Changelog{
					Title:    stringP("Changelog"),
					Releases: []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Title - Two Empty Trailing Lines": {
			Changelog: `# Changelog

`,
			Expected: expected{
				Result: &changelog.Changelog{
					Title:    stringP("Changelog"),
					Releases: []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Description - One Line": {
			Changelog: "Description",
			Expected: expected{
				Result: &changelog.Changelog{
					Description: stringP("Description"),
					Releases:    []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Description - One Line with Empty Trailing Line": {
			Changelog: `Description
`,
			Expected: expected{
				Result: &changelog.Changelog{
					Description: stringP("Description"),
					Releases:    []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Description - One Line with Two Empty Trailing Lines": {
			Changelog: `Description

`,
			Expected: expected{
				Result: &changelog.Changelog{
					Description: stringP("Description"),
					Releases:    []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Description - Multiple Lines": {
			Changelog: `Line 1
Line 2
Line 3`,
			Expected: expected{
				Result: &changelog.Changelog{
					Description: stringP("Line 1\nLine 2\nLine 3"),
					Releases:    []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Unreleased - Empty": {
			Changelog: `## [Unreleased]`,
			Expected: expected{
				Result: &changelog.Changelog{
					Unreleased: &changelog.Release{},
					Releases:   []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Unreleased - Link Definition": {
			Changelog: `## [Unreleased]

[Unreleased]: https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...HEAD`,
			Expected: expected{
				Result: &changelog.Changelog{
					Unreleased: &changelog.Release{
						URL: stringP("https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...HEAD"),
					},
					Releases: []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Unreleased - Inline Link": {
			Changelog: `## [Unreleased](https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...HEAD)
`,
			Expected: expected{
				Result: &changelog.Changelog{
					Unreleased: &changelog.Release{
						URL: stringP("https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...HEAD"),
					},
					Releases: []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Unreleased - Added Scope": {
			Changelog: `## [Unreleased]
### Added
- Change 1
- Change 2`,
			Expected: expected{
				Result: &changelog.Changelog{
					Unreleased: &changelog.Release{
						Changes: &changelog.Changes{
							Added: sliceOfStringsP([]string{
								"Change 1",
								"Change 2",
							}),
						},
					},
					Releases: []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Unreleased - Changed Scope": {
			Changelog: `## [Unreleased]
### Changed
- Change 1
- Change 2`,
			Expected: expected{
				Result: &changelog.Changelog{
					Unreleased: &changelog.Release{
						Changes: &changelog.Changes{
							Changed: sliceOfStringsP([]string{
								"Change 1",
								"Change 2",
							}),
						},
					},
					Releases: []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Unreleased - Deprecated Scope": {
			Changelog: `## [Unreleased]
### Deprecated
- Change 1
- Change 2`,
			Expected: expected{
				Result: &changelog.Changelog{
					Unreleased: &changelog.Release{
						Changes: &changelog.Changes{
							Deprecated: sliceOfStringsP([]string{
								"Change 1",
								"Change 2",
							}),
						},
					},
					Releases: []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Unreleased - Removed Scope": {
			Changelog: `## [Unreleased]
### Removed
- Change 1
- Change 2`,
			Expected: expected{
				Result: &changelog.Changelog{
					Unreleased: &changelog.Release{
						Changes: &changelog.Changes{
							Removed: sliceOfStringsP([]string{
								"Change 1",
								"Change 2",
							}),
						},
					},
					Releases: []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Unreleased - Fixed Scope": {
			Changelog: `## [Unreleased]
### Fixed
- Change 1
- Change 2`,
			Expected: expected{
				Result: &changelog.Changelog{
					Unreleased: &changelog.Release{
						Changes: &changelog.Changes{
							Fixed: sliceOfStringsP([]string{
								"Change 1",
								"Change 2",
							}),
						},
					},
					Releases: []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Unreleased - Security Scope": {
			Changelog: `## [Unreleased]
### Security
- Change 1
- Change 2`,
			Expected: expected{
				Result: &changelog.Changelog{
					Unreleased: &changelog.Release{
						Changes: &changelog.Changes{
							Security: sliceOfStringsP([]string{
								"Change 1",
								"Change 2",
							}),
						},
					},
					Releases: []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Unreleased - Notice Scope": {
			Changelog: `## [Unreleased]
Notice`,
			Expected: expected{
				Result: &changelog.Changelog{
					Unreleased: &changelog.Release{
						Changes: &changelog.Changes{
							Notice: stringP("Notice"),
						},
					},
					Releases: []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Unreleased - Full": {
			Changelog: `## [Unreleased]
Notice

### Added
- Change 1
- Change 2

### Changed
- Change 3
- Change 4

### Deprecated
- Change 5
- Change 6

### Removed
- Change 7
- Change 8

### Fixed
- Change 9
- Change 10

### Security
- Change 11
- Change 12`,
			Expected: expected{
				Result: &changelog.Changelog{
					Unreleased: &changelog.Release{
						Changes: &changelog.Changes{
							Notice: stringP("Notice"),
							Added: sliceOfStringsP([]string{
								"Change 1",
								"Change 2",
							}),
							Changed: sliceOfStringsP([]string{
								"Change 3",
								"Change 4",
							}),
							Deprecated: sliceOfStringsP([]string{
								"Change 5",
								"Change 6",
							}),
							Removed: sliceOfStringsP([]string{
								"Change 7",
								"Change 8",
							}),
							Fixed: sliceOfStringsP([]string{
								"Change 9",
								"Change 10",
							}),
							Security: sliceOfStringsP([]string{
								"Change 11",
								"Change 12",
							}),
						},
					},
					Releases: []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Unreleased - Complex Multiline Change": {
			Changelog: `## [Unreleased]
### Added
- Change 1:

` + "```yaml" + `
this:
  that:
  - A
  - B
` + "```" + `

- Change 2`,
			Expected: expected{
				Result: &changelog.Changelog{
					Unreleased: &changelog.Release{
						Changes: &changelog.Changes{
							Added: sliceOfStringsP([]string{
								"Change 1:\n\n```yaml\nthis:\n  that:\n  - A\n  - B\n```",
								"Change 2",
							}),
						},
					},
					Releases: []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Release - Empty": {
			Changelog: `## [0.0.1] - 2021-05-19`,
			Expected: expected{
				Result: &changelog.Changelog{
					Releases: []*changelog.Release{
						{
							Version: stringP("0.0.1"),
							Date:    dateP("2021-05-19"),
						},
					},
				},
				Error: "",
			},
		},
		"Release - Link Definition": {
			Changelog: `## [0.0.1] - 2021-05-19

[0.0.1]: https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1`,
			Expected: expected{
				Result: &changelog.Changelog{
					Releases: []*changelog.Release{
						{
							Version: stringP("0.0.1"),
							URL:     stringP("https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1"),
							Date:    dateP("2021-05-19"),
						},
					},
				},
				Error: "",
			},
		},
		"Release - Inline Link": {
			Changelog: `## [0.0.1](https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1) - 2021-05-19`,
			Expected: expected{
				Result: &changelog.Changelog{
					Releases: []*changelog.Release{
						{
							Version: stringP("0.0.1"),
							URL:     stringP("https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1"),
							Date:    dateP("2021-05-19"),
						},
					},
				},
				Error: "",
			},
		},
		"Release - Full": {
			Changelog: `## [0.0.1] - 2021-05-19
Notice

### Added
- Change 1
- Change 2
			
### Changed
- Change 3
- Change 4
			
### Deprecated
- Change 5
- Change 6
			
### Removed
- Change 7
- Change 8
			
### Fixed
- Change 9
- Change 10
			
### Security
- Change 11
- Change 12`,
			Expected: expected{
				Result: &changelog.Changelog{
					Releases: []*changelog.Release{
						{
							Version: stringP("0.0.1"),
							Date:    dateP("2021-05-19"),
							Changes: &changelog.Changes{
								Notice: stringP("Notice"),
								Added: sliceOfStringsP([]string{
									"Change 1",
									"Change 2",
								}),
								Changed: sliceOfStringsP([]string{
									"Change 3",
									"Change 4",
								}),
								Deprecated: sliceOfStringsP([]string{
									"Change 5",
									"Change 6",
								}),
								Removed: sliceOfStringsP([]string{
									"Change 7",
									"Change 8",
								}),
								Fixed: sliceOfStringsP([]string{
									"Change 9",
									"Change 10",
								}),
								Security: sliceOfStringsP([]string{
									"Change 11",
									"Change 12",
								}),
							},
						},
					},
				},
				Error: "",
			},
		},
		"Changelog - Title and Description": {
			Changelog: `# Changes
Notice`,
			Expected: expected{
				Result: &changelog.Changelog{
					Title:       stringP("Changes"),
					Description: stringP("Notice"),
					Releases:    []*changelog.Release{},
				},
				Error: "",
			},
		},
		"Changelog - Single Release and Description": {
			Changelog: `Notice

## [0.0.1] - 2021-05-19
### Added
- Change 1
- Change 2

[0.0.1]: https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1`,
			Expected: expected{
				Result: &changelog.Changelog{
					Description: stringP("Notice"),
					Releases: []*changelog.Release{
						{
							Version: stringP("0.0.1"),
							URL:     stringP("https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1"),
							Date:    dateP("2021-05-19"),
							Changes: &changelog.Changes{
								Added: sliceOfStringsP([]string{
									"Change 1",
									"Change 2",
								}),
							},
						},
					},
				},
				Error: "",
			},
		},
		"Changelog - Complex with Title, Description and Link Definitions": {
			Changelog: `# Changelog

Notice A
Notice B

## [0.0.2] - 2021-05-22

Notice

### Added
- Change 1
- Change 2
			
### Changed
- Change 3
- Change 4
			
### Deprecated
- Change 5
- Change 6
			
### Removed
- Change 7
- Change 8
			
### Fixed
- Change 9
- Change 10
			
### Security
- Change 11
- Change 12

## [0.0.1] - 2021-05-19
Notice
### Added
- Change 1
- Change 2
### Changed
- Change 3
- Change 4
### Deprecated
- Change 5
- Change 6
### Removed
- Change 7
- Change 8
### Fixed
- Change 9
- Change 10
### Security
- Change 11
- Change 12

[0.0.2]: https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1`,
			Expected: expected{
				Result: &changelog.Changelog{
					Title:       stringP("Changelog"),
					Description: stringP("Notice A\nNotice B"),
					Releases: []*changelog.Release{
						{
							Version: stringP("0.0.2"),
							URL:     stringP("https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...v0.0.2"),
							Date:    dateP("2021-05-22"),
							Changes: &changelog.Changes{
								Notice: stringP("Notice"),
								Added: sliceOfStringsP([]string{
									"Change 1",
									"Change 2",
								}),
								Changed: sliceOfStringsP([]string{
									"Change 3",
									"Change 4",
								}),
								Deprecated: sliceOfStringsP([]string{
									"Change 5",
									"Change 6",
								}),
								Removed: sliceOfStringsP([]string{
									"Change 7",
									"Change 8",
								}),
								Fixed: sliceOfStringsP([]string{
									"Change 9",
									"Change 10",
								}),
								Security: sliceOfStringsP([]string{
									"Change 11",
									"Change 12",
								}),
							},
						},
						{
							Version: stringP("0.0.1"),
							URL:     stringP("https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1"),
							Date:    dateP("2021-05-19"),
							Changes: &changelog.Changes{
								Notice: stringP("Notice"),
								Added: sliceOfStringsP([]string{
									"Change 1",
									"Change 2",
								}),
								Changed: sliceOfStringsP([]string{
									"Change 3",
									"Change 4",
								}),
								Deprecated: sliceOfStringsP([]string{
									"Change 5",
									"Change 6",
								}),
								Removed: sliceOfStringsP([]string{
									"Change 7",
									"Change 8",
								}),
								Fixed: sliceOfStringsP([]string{
									"Change 9",
									"Change 10",
								}),
								Security: sliceOfStringsP([]string{
									"Change 11",
									"Change 12",
								}),
							},
						},
					},
				},
				Error: "",
			},
		},
		"Changelog - Complex and Inconsistent": {
			Changelog: `Notice A
Notice B


## [Unreleased](https://github.com/anton-yurchenko/go-changelog/compare/v0.0.2...HEAD)
Notice

### Added
- Change 1:
A
B
- Change 2

## [0.0.2] - 2021-05-22

Notice

### Added
- Change 1
- Change 2
			
### Changed
- Change 3
- Change 4

## [0.0.1](https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1) - 2021-05-19
Notice
### Added
- Change 1
- Change 2
### Security
- Change 11
- Change 12

[0.0.2]: https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...v0.0.2`,
			Expected: expected{
				Result: &changelog.Changelog{
					Description: stringP("Notice A\nNotice B"),
					Unreleased: &changelog.Release{
						URL: stringP("https://github.com/anton-yurchenko/go-changelog/compare/v0.0.2...HEAD"),
						Changes: &changelog.Changes{
							Notice: stringP("Notice"),
							Added: sliceOfStringsP([]string{
								"Change 1:\nA\nB",
								"Change 2",
							}),
						},
					},
					Releases: []*changelog.Release{
						{
							Version: stringP("0.0.2"),
							URL:     stringP("https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...v0.0.2"),
							Date:    dateP("2021-05-22"),
							Changes: &changelog.Changes{
								Notice: stringP("Notice"),
								Added: sliceOfStringsP([]string{
									"Change 1",
									"Change 2",
								}),
								Changed: sliceOfStringsP([]string{
									"Change 3",
									"Change 4",
								}),
							},
						},
						{
							Version: stringP("0.0.1"),
							URL:     stringP("https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1"),
							Date:    dateP("2021-05-19"),
							Changes: &changelog.Changes{
								Notice: stringP("Notice"),
								Added: sliceOfStringsP([]string{
									"Change 1",
									"Change 2",
								}),
								Security: sliceOfStringsP([]string{
									"Change 11",
									"Change 12",
								}),
							},
						},
					},
				},
				Error: "",
			},
		},
		"Changelog - Yanked Releases": {
			Changelog: `## [0.0.2] - 2021-05-22 [YANKED]

Notice

### Added
- Change 1
- Change 2
			
### Changed
- Change 3
- Change 4

## [0.0.1](https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1) - 2021-05-19 [YANKED]
Notice
### Added
- Change 1
- Change 2
### Security
- Change 11
- Change 12

[0.0.2]: https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...v0.0.2`,
			Expected: expected{
				Result: &changelog.Changelog{
					Releases: []*changelog.Release{
						{
							Version: stringP("0.0.2"),
							URL:     stringP("https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...v0.0.2"),
							Date:    dateP("2021-05-22"),
							Yanked:  true,
							Changes: &changelog.Changes{
								Notice: stringP("Notice"),
								Added: sliceOfStringsP([]string{
									"Change 1",
									"Change 2",
								}),
								Changed: sliceOfStringsP([]string{
									"Change 3",
									"Change 4",
								}),
							},
						},
						{
							Version: stringP("0.0.1"),
							URL:     stringP("https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1"),
							Date:    dateP("2021-05-19"),
							Yanked:  true,
							Changes: &changelog.Changes{
								Notice: stringP("Notice"),
								Added: sliceOfStringsP([]string{
									"Change 1",
									"Change 2",
								}),
								Security: sliceOfStringsP([]string{
									"Change 11",
									"Change 12",
								}),
							},
						},
					},
				},
				Error: "",
			},
		},
		"Changelog - Empty Scopes": {
			Changelog: `Notice

## [0.0.1] - 2021-05-19
### Changed

### Fixed
### Added
- Change 1
- Change 2

[0.0.1]: https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1`,
			Expected: expected{
				Result: &changelog.Changelog{
					Description: stringP("Notice"),
					Releases: []*changelog.Release{
						{
							Version: stringP("0.0.1"),
							URL:     stringP("https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1"),
							Date:    dateP("2021-05-19"),
							Changes: &changelog.Changes{
								Added: sliceOfStringsP([]string{
									"Change 1",
									"Change 2",
								}),
							},
						},
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

		if name == "Changelog - Empty Scopes" {
			t.Log("got it!")
		}
		// prepare test case
		var p *changelog.Parser
		if name == "Filesystem Error" {
			m := new(mocks.Filesystem)
			m.On("Stat", "CHANGELOG.md").Return(nil, nil).Once()
			m.On("Open", "CHANGELOG.md").Return(nil, errors.New("reason")).Once()

			var err error
			p, err = changelog.NewParserWithFilesystem(m, "CHANGELOG.md")
			if err != nil {
				t.Errorf("error preparing test case: error creating parser: %v", err)
				continue
			}
		} else {
			err := afero.WriteFile(fs, "CHANGELOG.md", []byte(test.Changelog), 0644)
			if err != nil {
				t.Error("error preparing test case: error creating file CHANGELOG.md")
				continue
			}

			p, err = changelog.NewParserWithFilesystem(fs, "CHANGELOG.md")
			if err != nil {
				t.Errorf("error preparing test case: error creating parser: %v", err)
				continue
			}
		}

		// test
		c, err := p.Parse()
		a.Equal(test.Expected.Result, c)
		if test.Expected.Error != "" || err != nil {
			a.EqualError(err, test.Expected.Error)
		}

		// cleanup
		if name != "Filesystem Error" {
			if err := fs.Remove("CHANGELOG.md"); err != nil {
				t.Errorf("error cleanup: error removing file CHANGELOG.md: %v", err)
			}
		}
	}
}
