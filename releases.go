package changelog

import (
	"fmt"
	"regexp"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
)

// Releases is a slice of releases
type Releases []*Release

// Len return a total amount of Releases.
func (r Releases) Len() int {
	return len(r)
}

// Less compares versions of two releases.
func (r Releases) Less(i, j int) bool {
	return semver.Compare(fmt.Sprintf("v%v", *r[i].Version), fmt.Sprintf("v%v", *r[j].Version)) == -1
}

// Swap replaces positions of two releases.
func (r Releases) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// GetRelease returns a release for a provided version.
func (r Releases) GetRelease(version string) *Release {
	for _, release := range r {
		if *release.Version == version {
			return release
		}
	}

	return nil
}

// CreateRelease creates new empty release.
func (r *Releases) CreateRelease(version, date string) (*Release, error) {
	for _, e := range *r {
		if *e.Version == version {
			return nil, errors.New(fmt.Sprintf("version %v already exists", version))
		}
	}

	var d *time.Time
	dateMatcher := regexp.MustCompile(DateRegex)
	if dateMatcher.MatchString(date) {
		d = parseDate(date)
		if d == nil {
			return nil, errors.New(fmt.Sprintf("invalid date %v, expected format %v", date, DateFormat))
		}
	} else {
		return nil, errors.New(fmt.Sprintf("invalid date %v, expected to match regex %v", date, DateRegex))
	}

	var v string
	versionMatcher := regexp.MustCompile(SemVerRegex)
	if versionMatcher.MatchString(version) {
		v = version
	} else {
		return nil, errors.New(fmt.Sprintf("invalid semantic version %v, expected to match regex %v", version, SemVerRegex))
	}

	release := &Release{
		Changes: &Changes{},
		Date:    d,
		Version: &v,
	}

	*r = append(*r, release)

	return release, nil
}

// CreateReleaseWithURL creates new empty release.
//
// Identical to CreateRelease but with an extra step of adding a URL to the release.
func (r *Releases) CreateReleaseWithURL(version, date, url string) (*Release, error) {
	release, err := r.CreateRelease(version, date)
	if err != nil {
		return release, err
	}

	if err := release.SetURL(url); err != nil {
		return nil, err
	}

	return release, nil
}
