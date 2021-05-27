package changelog

const (
	// General
	EmptyLineRegex string = `^\s*$`
	URLRegex       string = `(([\w]+:)?\/\/)?(([\d\w]|%[a-fA-f\d]{2,2})+(:([\d\w]|%[a-fA-f\d]{2,2})+)?@)?([\d\w][-\d\w]{0,253}[\d\w]\.)+[\w]{2,4}(:[\d]+)?(\/([-+_~.\d\w]|%[a-fA-f\d]{2,2})*)*(\?(&amp;?([-+_~.\d\w]|%[a-fA-f\d]{2,2})=?)*)?(#([-+_~.\d\w]|%[a-fA-f\d]{2,2})*)?`
	SemVerRegex    string = `(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?`
	DateRegex      string = `([1-2][0-9][0-9][0-9])-([1-9]|[0][1-9]|[1][0-2])-([1-9]|[0][1-9]|[1-2][0-9]?|[3][0-1]?)`
	DateFormat     string = `2006-01-02`
	// Margins
	TitleRegex                       string = `^#\s*(?P<title>\S*)\s*$`
	UnreleasedTitleRegex             string = `^## \[(?P<title>Unreleased)\]$`
	UnreleasedTitleWithLinkRegex     string = `^## \[(?P<title>Unreleased)\]\((?P<url>` + URLRegex + `)\)$`
	VersionTitleRegex                string = `^## \[(?P<version>` + SemVerRegex + `)\] - (?P<date>` + DateRegex + `)(?P<yanked> \[YANKED\])?$`
	VersionTitleWithLinkRegex        string = `^## \[(?P<version>` + SemVerRegex + `)\]\((?P<url>` + URLRegex + `)\) - (?P<date>` + DateRegex + `)(?P<yanked> \[YANKED\])?$`
	MarkdownUnreleasedTitleLinkRegex string = `^\[(?P<title>Unreleased)\]: (?P<url>` + URLRegex + `)$`
	MarkdownVersionTitleLinkRegex    string = `^\[(?P<version>` + SemVerRegex + `)\]: (?P<url>` + URLRegex + `)$`
	// Scopes
	AddedScopeRegex      string = `^### (?P<scope>Added)$`
	ChangedScopeRegex    string = `^### (?P<scope>Changed)$`
	DeprecatedScopeRegex string = `^### (?P<scope>Deprecated)$`
	RemovedScopeRegex    string = `^### (?P<scope>Removed)$`
	FixedScopeRegex      string = `^### (?P<scope>Fixed)$`
	SecurityScopeRegex   string = `^### (?P<scope>Security)$`
	EntryRegex           string = `^(?P<marker>[-*+]\s*)(?P<entry>.*)$`
)
