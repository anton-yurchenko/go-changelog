package changelog

import (
	"fmt"
	"sort"
	"strings"
)

// Changelog reflects content of a complete changelog file
type Changelog struct {
	Title       *string
	Description *string
	Unreleased  *Release
	Releases    Releases
}

// ToString returns a Markdown formatted Changelog struct
func (c *Changelog) ToString() string {
	var o []string
	var defs []string

	if c.Title != nil {
		o = append(o, fmt.Sprintf("# %v\n", *c.Title))
	}

	if c.Description != nil {
		o = append(o, fmt.Sprintf("%v\n", *c.Description))
	}

	if c.Unreleased != nil {
		u, d := c.Unreleased.ToString()
		o = append(o, u)
		defs = append(defs, d)
	}

	sort.Sort(sort.Reverse(c.Releases))

	for _, release := range c.Releases {
		r, d := release.ToString()
		o = append(o, r)
		defs = append(defs, d)
	}

	o = append(o, defs...)

	return strings.Join(o, "\n")
}
