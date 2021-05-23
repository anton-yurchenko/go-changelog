package changelog

import (
	"fmt"
	"strings"
	"time"
)

// Release is a single version
type Release struct {
	Version *string
	Date    *time.Time
	URL     *string
	Changes *Changes
}

// ToString returns a Markdown formatted Release struct
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
