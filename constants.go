package changelog

const (
	// General
	EmptyLineRegex string = `^\s*$`
	// URLRegex taken from https://www.zimplicit.se/en/knowledge/validate-url-regular-expression
	URLRegex    string = `(([\w]+:)?\/\/)?(([\d\w]|%[a-fA-f\d]{2,2})+(:([\d\w]|%[a-fA-f\d]{2,2})+)?@)?([\d\w][-\d\w]{0,253}[\d\w]\.)+[\w]{2,4}(:[\d]+)?(\/([-+_~.\d\w]|%[a-fA-f\d]{2,2})*)*(\?(&amp;?([-+_~.\d\w]|%[a-fA-f\d]{2,2})=?)*)?(#([-+_~.\d\w]|%[a-fA-f\d]{2,2})*)?` // TODO: improve this expression
	SemVerRegex string = `(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?`
	DateRegex   string = `[0-9]{4}-[0-9]{2}-[0-9]{2}` // TODO: can be improbed
	// Margins
	TitleRegex                       string = `^#\s*(?P<title>\S*)\s*$`
	UnreleasedTitleRegex             string = `^## \[(?P<title>Unreleased)\]$`
	UnreleasedTitleWithLinkRegex     string = `^## \[(?P<title>Unreleased)\]\((?P<url>` + URLRegex + `)\)$`
	VersionTitleRegex                string = `^## \[(?P<version>` + SemVerRegex + `)\] - ` + DateRegex + `$`
	VersionTitleWithLinkRegex        string = `^## \[(?P<version>` + SemVerRegex + `)\]\((?P<url>` + URLRegex + `)\) - ` + DateRegex + `$`
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
