package changelog

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Parser is basically a runtime that holds a raw changelog content,
// key file margins, filesystem backend and other attributes.
type Parser struct {
	Filepath   string
	filesystem afero.Fs
	Buffer     []string
	margins    struct {
		lines      []int
		title      *int
		unreleased *int
		releases   []int
		links      []int
		added      []int
		changed    []int
		deprecated []int
		removed    []int
		fixed      []int
		security   []int
	}
}

// NewParser creates a new Changelog Parser.
//
// *`filepath` is validated for readability.*
func NewParser(filepath string) (*Parser, error) {
	return NewParserWithFilesystem(afero.NewOsFs(), filepath)
}

// NewParser creates a new Changelog Parser using non default (OS) filesystem.
// Possible options for `filesystems` are: [`afero.NewOsFs()`,`afero.NewMemMapFs()`]
//
// *`filepath` is validated*
func NewParserWithFilesystem(filesystem afero.Fs, filepath string) (*Parser, error) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) || os.IsPermission(err) {
		return nil, errors.Wrapf(err, "error accessing a file %v", filepath)
	}

	return &Parser{
		Filepath:   filepath,
		filesystem: filesystem}, nil
}

// Parse a changelog file and return a Changelog struct
func (p *Parser) Parse() (*Changelog, error) {
	o := new(Changelog)

	if err := p.loadBuffer(); err != nil {
		return nil, errors.Wrap(err, "error loading a buffer")
	}

	p.identifyMargins()
	o.Title = p.parseTitle()
	o.Description = p.parseDescription()
	o.Unreleased = p.parseUnreleased()
	o.Releases = p.parseReleases()

	return o, nil
}

func (p *Parser) loadBuffer() error {
	lines := make([]string, 0)

	file, err := p.filesystem.Open(p.Filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	p.Buffer = lines
	return nil
}

func (p *Parser) identifyMargins() {
	matchers := map[string]*regexp.Regexp{
		"Title":                       regexp.MustCompile(TitleRegex),
		"UnreleasedTitle":             regexp.MustCompile(UnreleasedTitleRegex),
		"UnreleasedTitleWithLink":     regexp.MustCompile(UnreleasedTitleWithLinkRegex),
		"VersionTitle":                regexp.MustCompile(VersionTitleRegex),
		"VersionTitleWithLink":        regexp.MustCompile(VersionTitleWithLinkRegex),
		"MarkdownUnreleasedTitleLink": regexp.MustCompile(MarkdownUnreleasedTitleLinkRegex),
		"MarkdownVersionTitleLink":    regexp.MustCompile(MarkdownVersionTitleLinkRegex),
		"AddedScope":                  regexp.MustCompile(AddedScopeRegex),
		"ChangedScope":                regexp.MustCompile(ChangedScopeRegex),
		"DeprecatedScope":             regexp.MustCompile(DeprecatedScopeRegex),
		"RemovedScope":                regexp.MustCompile(RemovedScopeRegex),
		"FixedScope":                  regexp.MustCompile(FixedScopeRegex),
		"SecurityScope":               regexp.MustCompile(SecurityScopeRegex),
	}

	for i, l := range p.Buffer {
		for k, m := range matchers {
			if m.MatchString(l) {
				switch k {
				case "Title":
					n := i
					p.margins.title = &n
				case "UnreleasedTitle":
					n := i
					p.margins.unreleased = &n
				case "UnreleasedTitleWithLink":
					n := i
					p.margins.unreleased = &n
				case "VersionTitle":
					p.margins.releases = append(p.margins.releases, i)
				case "VersionTitleWithLink":
					p.margins.releases = append(p.margins.releases, i)
				case "MarkdownUnreleasedTitleLink":
					p.margins.links = append(p.margins.links, i)
				case "MarkdownVersionTitleLink":
					p.margins.links = append(p.margins.links, i)
				case "AddedScope":
					p.margins.added = append(p.margins.added, i)
				case "ChangedScope":
					p.margins.changed = append(p.margins.changed, i)
				case "DeprecatedScope":
					p.margins.deprecated = append(p.margins.deprecated, i)
				case "RemovedScope":
					p.margins.removed = append(p.margins.removed, i)
				case "FixedScope":
					p.margins.fixed = append(p.margins.fixed, i)
				case "SecurityScope":
					p.margins.security = append(p.margins.security, i)
				}

				p.margins.lines = append(p.margins.lines, i)
			}
		}
	}
}

func (p *Parser) parseTitle() *string {
	if p.margins.title != nil {
		m := regexp.MustCompile(TitleRegex)
		x := m.ReplaceAllString(p.Buffer[*p.margins.title], "${1}")
		return &x
	}

	return nil
}

func (p *Parser) parseDescription() *string {
	if p.margins.title != nil {
		n := p.getNextMarginLine(*p.margins.title)
		if n != nil {
			x := strings.Join(trimLeadingAndTrailingEmptyLines(p.Buffer[*p.margins.title+1:*n]), "\n")
			return &x
		}
	} else if p.margins.lines[0] > 0 {
		n := p.getNextMarginLine(0)
		if n == nil {
			x := strings.Join(trimLeadingAndTrailingEmptyLines(p.Buffer[0:p.margins.lines[0]-1]), "\n")
			return &x
		}
	}

	return nil
}

func trimLeadingAndTrailingEmptyLines(content []string) []string {
	newContent := make([]string, 0)
	started := false
	trailing := 0
	m := regexp.MustCompile(EmptyLineRegex)

	for _, l := range content {
		if m.MatchString(l) {
			if !started {
				continue
			}

			trailing++
			newContent = append(newContent, l)
		} else {
			started = true
			trailing = 0
			newContent = append(newContent, l)
		}
	}

	return newContent[:len(newContent)-trailing]
}

func (p *Parser) getNextMarginLine(current int) *int {
	for i, n := range p.margins.lines {
		if current == n {
			if i == len(p.margins.lines)-1 {
				return &current
			}

			next := p.margins.lines[i+1]
			return &next
		}
	}

	return nil
}

func (p *Parser) parseUnreleased() *Release {
	if p.margins.unreleased != nil {
		// NOTE: `parseRelease` function will try to parse from an inline link and fall back to definitions parsing
		return p.parseRelease(nil, *p.margins.unreleased, UnreleasedTitleWithLinkRegex)
	}

	return nil
}

func (p *Parser) parseReleases() []*Release {
	releases := make([]*Release, 0)

	matchers := []*regexp.Regexp{
		regexp.MustCompile(VersionTitleRegex),
		regexp.MustCompile(VersionTitleWithLinkRegex),
	}

	for _, n := range p.margins.releases {
		for _, m := range matchers {
			/* NOTE: match inline link or definition, but always pass an inline link regex.
			`parseRelease` function will try to parse from it and fall back to definitions parsing
			*/
			if m.MatchString(p.Buffer[n]) {
				v := m.ReplaceAllString(p.Buffer[n], "${1}")
				releases = append(releases, p.parseRelease(&v, n, VersionTitleWithLinkRegex))
			}
		}
	}

	return releases
}

func (p *Parser) parseRelease(version *string, startingLine int, titleWithLinkRegex string) *Release {
	release := new(Release)

	/* NOTE: parse URL
	Try to parse from title `## [Unreleased]()`/`## [<version>](<url>) - <date>`,
	otherwise, try to parse the URL from definitions
	*/
	m := regexp.MustCompile(titleWithLinkRegex)
	if m.MatchString(p.Buffer[startingLine]) {
		x := m.ReplaceAllString(p.Buffer[startingLine], "${7}")
		release.URL = &x
	} else {
		release.URL = p.parseLinkURL(version)
	}

	// NOTE: parse changes
	n := p.getReleaseEndLine(startingLine)
	if startingLine == n {
		release.Changes = nil
	} else {
		release.Changes = p.parseChanges(startingLine+1, n)
	}

	return release
}

func (p *Parser) parseLinkURL(version *string) *string {
	var matcher *regexp.Regexp
	var position int

	if version == nil {
		position = 2
		matcher = regexp.MustCompile(MarkdownUnreleasedTitleLinkRegex)
	} else {
		position = 7
		matcher = regexp.MustCompile(MarkdownVersionTitleLinkRegex)
	}

	for _, n := range p.margins.links {
		if matcher.MatchString(p.Buffer[n]) {
			if version == nil {
				x := matcher.ReplaceAllString(p.Buffer[n], fmt.Sprintf("${%v}", position))
				return &x
			} else if version != nil {
				if matcher.ReplaceAllString(p.Buffer[n], "${1}") == *version {
					x := matcher.ReplaceAllString(p.Buffer[n], fmt.Sprintf("${%v}", position))
					return &x
				}
			}
		}
	}

	return nil
}

func (p *Parser) getReleaseEndLine(startingLine int) int {
	sections := make([]int, 0)
	if p.margins.title != nil {
		sections = append(sections, *p.margins.title)
	}
	if p.margins.unreleased != nil {
		sections = append(sections, *p.margins.unreleased)
	}
	sections = append(sections, p.margins.releases...)
	sections = append(sections, p.margins.links...)

	sort.Ints(sections)

	i := getNextItem(startingLine, sections)
	if i != nil {
		return *i
	}

	return len(p.Buffer) - 1
}

func getNextItem(current int, array []int) *int {
	for i, n := range array {
		if n == current {
			if i != len(array)-1 {
				x := array[i+1] - 1
				return &x
			}
		}
	}

	return nil
}

func (p *Parser) parseChanges(startingLine, endLine int) *Changes {
	changes := new(Changes)

	scopeLines := map[string]*int{
		"Added":      getScopeLine(startingLine, endLine, p.margins.added),
		"Changed":    getScopeLine(startingLine, endLine, p.margins.changed),
		"Deprecated": getScopeLine(startingLine, endLine, p.margins.deprecated),
		"Removed":    getScopeLine(startingLine, endLine, p.margins.removed),
		"Fixed":      getScopeLine(startingLine, endLine, p.margins.fixed),
		"Security":   getScopeLine(startingLine, endLine, p.margins.security),
	}

	lines := make([]int, 0)
	for _, n := range scopeLines {
		if n != nil {
			lines = append(lines, *n)
		}
	}
	sort.Ints(lines)

	if len(lines) > 0 {
		val := strings.Join(trimLeadingAndTrailingEmptyLines(p.Buffer[startingLine:lines[0]]), "\n")
		if val != "" {
			changes.Notice = &val
		}
	}

	notEmpty := false
	for name, start := range scopeLines {
		if start != nil {
			e := getNextItem(*start, lines)
			var end int
			if e == nil {
				end = endLine
			} else {
				end = *e
			}
			entries := p.parseScopeEntries(*start, end)

			switch name {
			case "Added":
				notEmpty = true
				changes.Added = entries
			case "Changed":
				notEmpty = true
				changes.Changed = entries
			case "Deprecated":
				notEmpty = true
				changes.Deprecated = entries
			case "Removed":
				notEmpty = true
				changes.Removed = entries
			case "Fixed":
				notEmpty = true
				changes.Fixed = entries
			case "Security":
				notEmpty = true
				changes.Security = entries
			}
		}
	}

	if notEmpty {
		return changes
	}
	return nil
}

func getScopeLine(startingLine, endLine int, lines []int) *int {
	var s *int
	counter := 0

	for _, l := range lines {
		if l >= startingLine && l <= endLine {
			counter++
			x := l
			s = &x
		}
	}

	if counter != 1 {
		return nil
	}

	return s
}

func (p *Parser) parseScopeEntries(startingLine, endLine int) *[]string {
	entries := make([]string, 0)
	entryLines := make([]int, 0)

	m := regexp.MustCompile(EntryRegex)
	for i := startingLine; i <= endLine; i++ {
		if m.MatchString(p.Buffer[i]) {
			entryLines = append(entryLines, i)
		}
	}

	for i, n := range entryLines {
		start := n
		var end int
		if i == len(entryLines)-1 {
			end = endLine
		} else {
			end = entryLines[i+1]
		}

		entry := []string{m.ReplaceAllString(p.Buffer[start], "${2}")}
		for ii := start + 1; ii < end; ii++ {
			entry = append(entry, p.Buffer[ii])
		}

		entries = append(entries, strings.Join(trimLeadingAndTrailingEmptyLines(entry), "\n"))
	}

	return &entries
}
