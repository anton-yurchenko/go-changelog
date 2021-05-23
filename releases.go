package changelog

import (
	"fmt"

	"golang.org/x/mod/semver"
)

// Releases is a slice of pointers to a Release
type Releases []*Release

// Len return a total amount of Releases
func (r Releases) Len() int {
	return len(r)
}

// Less compares versions of two releases
func (r Releases) Less(i, j int) bool {
	return semver.Compare(fmt.Sprintf("v%v", *r[i].Version), fmt.Sprintf("v%v", *r[j].Version)) == -1
}

// Swap places of two releases
func (r Releases) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
