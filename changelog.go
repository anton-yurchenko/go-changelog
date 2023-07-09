package changelog

import (
	"fmt"
	"io/fs"
	"net/url"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Filesystem is an interface of a filesystem.
//
// Possible options: [afero.NewOsFs(), afero.NewMemMapFs()].
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

// ToString returns a Markdown formatted Changelog struct.
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
//
// Possible options for Filesystem are: [afero.NewOsFs(), afero.NewMemMapFs()].
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

// NewChangelog returns an empty changelog.
func NewChangelog() *Changelog {
	c := new(Changelog)
	c.Releases = make(Releases, 0)

	return c
}

// SetTitle updates a title of the changelog.
func (c *Changelog) SetTitle(title string) {
	*c.Title = title
}

// SetDescription updates a description of the changelog.
func (c *Changelog) SetDescription(description string) {
	*c.Description = description
}

// SetUnreleasedURL configures a markdown URL for an Unreleased section.
func (c *Changelog) SetUnreleasedURL(link string) error {
	_, err := url.Parse(link)
	if err != nil {
		return err
	}

	if c.Unreleased == nil {
		c.Unreleased = &Release{
			URL: &link,
		}
	} else {
		c.Unreleased.URL = &link
	}

	return nil
}

// AddUnreleasedChange adds a scoped change to Unreleased section.
//
// Supported scopes: [added, changed, deprecated, removed, fixed, security].
func (c *Changelog) AddUnreleasedChange(scope string, change string) error {
	if c.Unreleased == nil {
		c.Unreleased = &Release{
			Changes: &Changes{},
		}
	}

	if c.Unreleased.Changes == nil {
		c.Unreleased.Changes = new(Changes)
	}

	return c.Unreleased.Changes.AddChange(scope, change)
}

// GetRelease returns a release for a provided version.
//
// This is a helper function that wraps Releases.GetRelease function.
func (c *Changelog) GetRelease(version string) *Release {
	return c.Releases.GetRelease(version)
}

// CreateReleaseFromUnreleased creates a new release with all the changes from Unreleased section.
// This will also cleanup the Unreleased section.
func (c *Changelog) CreateReleaseFromUnreleased(version, date string) (*Release, error) {
	if c.Unreleased == nil || c.Unreleased.Changes == nil {
		return nil, errors.New("missing 'Unreleased' section")
	}

	r, err := c.CreateRelease(version, date)
	if err != nil {
		return nil, err
	}

	r.Changes = c.Unreleased.Changes
	c.Unreleased.Changes = nil

	return r, nil
}

// CreateReleaseFromUnreleased creates a new release with all the changes from Unreleased section.
// This will also cleanup the Unreleased section.
// Identical to CreateReleaseFromUnreleased but with an extra step of adding a URL to the release.
func (c *Changelog) CreateReleaseFromUnreleasedWithURL(version, date, url string) (*Release, error) {
	r, err := c.CreateReleaseFromUnreleased(version, date)
	if err != nil {
		return nil, err
	}

	if err := r.SetURL(url); err != nil {
		return nil, err
	}

	return r, nil
}

// CreateRelease creates new empty release.
//
// This is a helper function that wraps Releases.CreateRelease function.
func (c *Changelog) CreateRelease(version, date string) (*Release, error) {
	return c.Releases.CreateRelease(version, date)
}

// CreateReleaseWithURL creates new empty release.
//
// This is a helper function that wraps Releases.CreateReleaseWithURL function.
//
// Identical to CreateRelease but with an extra step of adding a URL to the release.
func (c *Changelog) CreateReleaseWithURL(version, date, url string) (*Release, error) {
	return c.Releases.CreateReleaseWithURL(version, date, url)
}
