package changelog

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type Filesystem interface {
	Stat(string) (fs.FileInfo, error)
	Open(string) (afero.File, error)
	Create(string) (afero.File, error)
}

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

// SaveToFile formats the changelog struct according to a predefined format
// and prints it to file.
// Possible options for `filesystems` are: [`afero.NewOsFs()`,`afero.NewMemMapFs()`]
func (c *Changelog) SaveToFile(filesystem Filesystem, filepath string) error {
	f, err := filesystem.Create(filepath)
	if err != nil {
		return errors.Wrap(err, "error creating a file")
	}
	defer f.Close()

	_, err = f.WriteString(c.ToString())
	if err != nil {
		return errors.Wrap(err, "error writing to file")
	}

	if err := f.Sync(); err != nil {
		return errors.Wrap(err, "error committing file content to disk")
	}

	return nil
}

// NewChangelog returns an empty changelog
func NewChangelog() *Changelog {
	c := new(Changelog)
	c.Releases = make(Releases, 0)

	return c
}
