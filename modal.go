package changelog

// Changelog reflects content of a complete changelog file
type Changelog struct {
	Title       *string
	Description *string
	Unreleased  *Release
	Releases    []*Release
}

// Release is a single version
type Release struct {
	Version *string
	URL     *string
	Changes *Changes
}

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
