package changelog

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Changes are scoped changelog entries for a single version.
type Changes struct {
	Added      *[]string
	Changed    *[]string
	Deprecated *[]string
	Fixed      *[]string
	Notice     *string
	Removed    *[]string
	Security   *[]string
}

// ToString returns a Markdown formatted Changes struct.
func (c *Changes) ToString() string {
	var o []string
	if c.Notice != nil {
		o = append(o, fmt.Sprintf("%v\n", *c.Notice))
	}

	if c.Security != nil {
		o = append(o, "### Security\n", fmt.Sprintf("%v\n", scopeToString(c.Security)))
	}

	if c.Changed != nil {
		o = append(o, "### Changed\n", fmt.Sprintf("%v\n", scopeToString(c.Changed)))
	}

	if c.Added != nil {
		o = append(o, "### Added\n", fmt.Sprintf("%v\n", scopeToString(c.Added)))
	}

	if c.Removed != nil {
		o = append(o, "### Removed\n", fmt.Sprintf("%v\n", scopeToString(c.Removed)))
	}

	if c.Fixed != nil {
		o = append(o, "### Fixed\n", fmt.Sprintf("%v\n", scopeToString(c.Fixed)))
	}

	if c.Deprecated != nil {
		o = append(o, "### Deprecated\n", fmt.Sprintf("%v\n", scopeToString(c.Deprecated)))
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

// AddNotice adds a notice to the changes.
func (c *Changes) AddNotice(notice string) {
	*c.Notice = notice
}

// AddChange adds a scoped change.
//
// Supported scopes: [added, changed, deprecated, removed, fixed, security].
func (c *Changes) AddChange(scope string, change string) error {
	changesList := []string{change}

	switch strings.ToLower(scope) {
	case "added":
		if change == "" {
			return nil
		}

		if c.Added == nil {
			c.Added = &changesList
		} else {
			*c.Added = append(*c.Added, change)
		}
	case "changed":
		if change == "" {
			return nil
		}

		if c.Changed == nil {
			c.Changed = &changesList
		} else {
			*c.Changed = append(*c.Changed, change)
		}
	case "deprecated":
		if change == "" {
			return nil
		}

		if c.Deprecated == nil {
			c.Deprecated = &changesList
		} else {
			*c.Deprecated = append(*c.Deprecated, change)
		}
	case "removed":
		if change == "" {
			return nil
		}

		if c.Removed == nil {
			c.Removed = &changesList
		} else {
			*c.Removed = append(*c.Removed, change)
		}
	case "fixed":
		if change == "" {
			return nil
		}

		if c.Fixed == nil {
			c.Fixed = &changesList
		} else {
			*c.Fixed = append(*c.Fixed, change)
		}
	case "security":
		if change == "" {
			return nil
		}

		if c.Security == nil {
			c.Security = &changesList
		} else {
			*c.Security = append(*c.Security, change)
		}
	default:
		return errors.New(fmt.Sprintf("unexpected scope: %v (supported: [added,changed,deprecated,removed,fixed,security])", scope))
	}

	return nil
}
