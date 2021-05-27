package changelog

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
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

func (c *Changes) AddNotice(notice string) {
	*c.Notice = notice
}

func (c *Changes) AddChange(scope string, change string) error {
	changesList := []string{change}

	switch scope {
	case "Added":
		if change == "" {
			return nil
		}

		if c.Added == nil {
			c.Added = &changesList
		} else {
			*c.Added = append(*c.Added, change)
		}
	case "Changed":
		if change == "" {
			return nil
		}

		if c.Changed == nil {
			c.Changed = &changesList
		} else {
			*c.Changed = append(*c.Changed, change)
		}
	case "Deprecated":
		if change == "" {
			return nil
		}

		if c.Deprecated == nil {
			c.Deprecated = &changesList
		} else {
			*c.Deprecated = append(*c.Deprecated, change)
		}
	case "Removed":
		if change == "" {
			return nil
		}

		if c.Removed == nil {
			c.Removed = &changesList
		} else {
			*c.Removed = append(*c.Removed, change)
		}
	case "Fixed":
		if change == "" {
			return nil
		}

		if c.Fixed == nil {
			c.Fixed = &changesList
		} else {
			*c.Fixed = append(*c.Fixed, change)
		}
	case "Security":
		if change == "" {
			return nil
		}

		if c.Security == nil {
			c.Security = &changesList
		} else {
			*c.Security = append(*c.Security, change)
		}
	default:
		return errors.New(fmt.Sprintf("unexpected scope: %v (supported: [Added,Changed,Deprecated,Removed,Fixed,Security])", scope))
	}

	return nil
}
