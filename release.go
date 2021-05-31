package changelog

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Release is a single changelog version
type Release struct {
	Version *string
	Date    *time.Time
	Yanked  bool
	URL     *string
	Changes *Changes
}

// ToString returns a Markdown formatted Release struct.
func (r *Release) ToString() (string, string) {
	var o []string
	var u string

	if r.Version != nil {
		if r.Date != nil {
			o = append(o, fmt.Sprintf("## [%v] - %v", *r.Version, formatDate(r.Date)))
		} else {
			o = append(o, fmt.Sprintf("## [%v]", *r.Version))
		}
	} else {
		o = append(o, "## [Unreleased]")
	}

	if r.Changes != nil {
		o = append(o, r.Changes.ToString())
	}

	if r.URL != nil {
		if r.Version != nil {
			u = fmt.Sprintf("[%v]: %v", *r.Version, *r.URL)
		} else {
			u = fmt.Sprintf("[Unreleased]: %v", *r.URL)
		}
	}

	return strings.Join(o, "\n"), u
}

func formatDate(date *time.Time) string {
	return date.Format(DateFormat)
}

// SetVersion configures a Semantic Version of a release.
func (r *Release) SetVersion(version string) error {
	m := regexp.MustCompile(SemVerRegex)
	if m.MatchString(version) {
		r.Version = &version
		return nil
	}

	return errors.New(fmt.Sprintf("invalid semantic version %v, expected to match regex %v", version, SemVerRegex))
}

// SetDate configures a date of the release.
// Expected format: YYYY-MM-DD
func (r *Release) SetDate(date string) error {
	m := regexp.MustCompile(DateRegex)
	if m.MatchString(date) {
		d := parseDate(date)
		if d == nil {
			return errors.New(fmt.Sprintf("invalid date %v, expected format %v", date, DateFormat))
		}

		r.Date = d
		return nil
	}

	return errors.New(fmt.Sprintf("invalid date %v, expected to match regex %v", date, DateRegex))
}

// SetURL configures a URL of the release
func (r *Release) SetURL(link string) error {
	_, err := url.Parse(link)
	if err != nil {
		return err
	}

	r.URL = &link
	return nil
}

// AddNotice adds a notice to the release.
//
// This is a helper function that wraps Changes.AddNotice function.
func (r *Release) AddNotice(notice string) {
	r.Changes.AddNotice(notice)
}

// AddChange adds a scoped change to the release.
//
// This is a helper function that wraps Changes.AddChange function.
func (r *Release) AddChange(scope string, change string) error {
	return r.Changes.AddChange(scope, change)
}
