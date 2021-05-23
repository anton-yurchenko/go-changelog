package changelog

import (
	"fmt"
	"strings"
)

// Changes are scoped changelog entries for a single version.
type Changes struct {
	Notice     *string
	Added      *[]string
	Changed    *[]string
	Deprecated *[]string
	Removed    *[]string
	Fixed      *[]string
	Security   *[]string
}

// ToString returns a Markdown formatted Changes struct
func (c *Changes) ToString() string {
	var o []string
	if c.Notice != nil {
		o = append(o, fmt.Sprintf("%v\n", *c.Notice))
	}

	if c.Added != nil {
		o = append(o, "### Added", fmt.Sprintf("%v\n", scopeToString(c.Added)))
	}

	if c.Changed != nil {
		o = append(o, "### Changed", fmt.Sprintf("%v\n", scopeToString(c.Changed)))
	}

	if c.Deprecated != nil {
		o = append(o, "### Deprecated", fmt.Sprintf("%v\n", scopeToString(c.Deprecated)))
	}

	if c.Removed != nil {
		o = append(o, "### Removed", fmt.Sprintf("%v\n", scopeToString(c.Removed)))
	}

	if c.Fixed != nil {
		o = append(o, "### Fixed", fmt.Sprintf("%v\n", scopeToString(c.Fixed)))
	}

	if c.Security != nil {
		o = append(o, "### Security", fmt.Sprintf("%v\n", scopeToString(c.Security)))
	}

	return strings.Join(o, "\n")
}

func scopeToString(scope *[]string) string {
	var o []string
	for _, c := range *scope {
		o = append(o, fmt.Sprintf("- %v", c))
	}

	return strings.Join(o, "\n")
}
