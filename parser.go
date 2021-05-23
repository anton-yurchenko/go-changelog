package changelog

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Parser is basically a runtime that holds a raw changelog content,
// key file Margins, filesystem backend and other attributes.
type Parser struct {
	Filepath   string
	Filesystem Filesystem
	Buffer     []string
	Margins    struct {
		Lines      []int
		Title      *int
		Unreleased *int
		Releases   []int
		Links      []int
		Added      []int
		Changed    []int
		Deprecated []int
		Removed    []int
		Fixed      []int
		Security   []int
	}
}

type Filesystem interface {
	Stat(string) (fs.FileInfo, error)
	Open(string) (afero.File, error)
}

// NewParser creates a new Changelog Parser.
func NewParser(filepath string) (*Parser, error) {
	return NewParserWithFilesystem(afero.NewOsFs(), filepath)
}

// NewParser creates a new Changelog Parser using non default (OS) filesystem.
// Possible options for `filesystems` are: [`afero.NewOsFs()`,`afero.NewMemMapFs()`]
func NewParserWithFilesystem(filesystem Filesystem, filepath string) (*Parser, error) {
	_, err := filesystem.Stat(filepath)
	if os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "file %v not found", filepath)
	}

	return &Parser{
		Filepath:   filepath,
		Filesystem: filesystem}, nil
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

	file, err := p.Filesystem.Open(p.Filepath)
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
					p.Margins.Title = &n
				case "UnreleasedTitle":
					n := i
					p.Margins.Unreleased = &n
				case "UnreleasedTitleWithLink":
					n := i
					p.Margins.Unreleased = &n
				case "VersionTitle":
					p.Margins.Releases = append(p.Margins.Releases, i)
				case "VersionTitleWithLink":
					p.Margins.Releases = append(p.Margins.Releases, i)
				case "MarkdownUnreleasedTitleLink":
					p.Margins.Links = append(p.Margins.Links, i)
				case "MarkdownVersionTitleLink":
					p.Margins.Links = append(p.Margins.Links, i)
				case "AddedScope":
					p.Margins.Added = append(p.Margins.Added, i)
				case "ChangedScope":
					p.Margins.Changed = append(p.Margins.Changed, i)
				case "DeprecatedScope":
					p.Margins.Deprecated = append(p.Margins.Deprecated, i)
				case "RemovedScope":
					p.Margins.Removed = append(p.Margins.Removed, i)
				case "FixedScope":
					p.Margins.Fixed = append(p.Margins.Fixed, i)
				case "SecurityScope":
					p.Margins.Security = append(p.Margins.Security, i)
				}

				p.Margins.Lines = append(p.Margins.Lines, i)
			}
		}
	}
}

func (p *Parser) parseTitle() *string {
	if p.Margins.Title != nil {
		m := regexp.MustCompile(TitleRegex)
		x := m.ReplaceAllString(p.Buffer[*p.Margins.Title], "${1}")
		return &x
	}

	return nil
}

func (p *Parser) parseDescription() *string {
	var o string

	if p.Margins.Title != nil {
		s := *p.Margins.Title + 1
		var e *int

		if len(p.Margins.Lines) == 1 {
			t := len(p.Buffer)
			e = &t
		} else {
			e = p.getNextMarginLine(*p.Margins.Title)
		}

		o = strings.Join(trimLeadingAndTrailingEmptyLines(p.Buffer[s:*e]), "\n")
	} else {
		if len(p.Margins.Lines) != 0 {
			if p.Margins.Lines[0] != 0 {
				s := 0
				e := p.Margins.Lines[0] - 1

				o = strings.Join(trimLeadingAndTrailingEmptyLines(p.Buffer[s:e]), "\n")
			}
		} else {
			o = strings.Join(trimLeadingAndTrailingEmptyLines(p.Buffer), "\n")
		}
	}

	if o == "" {
		return nil
	}
	return &o
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
	var o int

	for i, n := range p.Margins.Lines {
		if current == n {
			o = p.Margins.Lines[i+1]
		}
	}

	return &o
}

func (p *Parser) parseUnreleased() *Release {
	if p.Margins.Unreleased != nil {
		// NOTE: `parseRelease` function will try to parse from an inline link and fall back to definitions parsing
		return p.parseRelease(nil, *p.Margins.Unreleased, UnreleasedTitleWithLinkRegex)
	}

	return nil
}

func (p *Parser) parseReleases() Releases {
	releases := make([]*Release, 0)

	matchers := []*regexp.Regexp{
		regexp.MustCompile(VersionTitleRegex),
		regexp.MustCompile(VersionTitleWithLinkRegex),
	}

	for _, n := range p.Margins.Releases {
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
	if version != nil {
		release.Version = version
	}

	/* NOTE: parse URL
	Try to parse from title `## [Unreleased](<url>)`/`## [<version>](<url>) - <date>`,
	otherwise, try to parse the URL from definitions
	*/
	m := regexp.MustCompile(titleWithLinkRegex)
	if m.MatchString(p.Buffer[startingLine]) {
		var x string
		if version == nil {
			// Unreleased
			x = m.ReplaceAllString(p.Buffer[startingLine], "${2}")
			release.Date = nil
		} else {
			// Release
			x = m.ReplaceAllString(p.Buffer[startingLine], "${7}")
		}

		release.URL = &x
	} else {
		release.URL = p.parseLinkURL(version)
	}

	// NOTE: parse date
	if version != nil {
		m1 := regexp.MustCompile(VersionTitleRegex)
		m2 := regexp.MustCompile(VersionTitleWithLinkRegex)

		if m1.MatchString(p.Buffer[startingLine]) {
			release.Date = parseDate(m1.ReplaceAllString(p.Buffer[startingLine], "${7}"))
		} else if m2.MatchString(p.Buffer[startingLine]) {
			release.Date = parseDate(m2.ReplaceAllString(p.Buffer[startingLine], "${23}"))
		}
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

func parseDate(date string) *time.Time {
	t, err := time.Parse(DateFormat, date)
	if err != nil {
		return nil
	}

	return &t
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

	for _, n := range p.Margins.Links {
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
	if p.Margins.Title != nil {
		sections = append(sections, *p.Margins.Title)
	}
	if p.Margins.Unreleased != nil {
		sections = append(sections, *p.Margins.Unreleased)
	}
	sections = append(sections, p.Margins.Releases...)
	sections = append(sections, p.Margins.Links...)

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
		"Added":      getScopeLine(startingLine, endLine, p.Margins.Added),
		"Changed":    getScopeLine(startingLine, endLine, p.Margins.Changed),
		"Deprecated": getScopeLine(startingLine, endLine, p.Margins.Deprecated),
		"Removed":    getScopeLine(startingLine, endLine, p.Margins.Removed),
		"Fixed":      getScopeLine(startingLine, endLine, p.Margins.Fixed),
		"Security":   getScopeLine(startingLine, endLine, p.Margins.Security),
	}

	lines := make([]int, 0)
	for _, n := range scopeLines {
		if n != nil {
			lines = append(lines, *n)
		}
	}
	sort.Ints(lines)

	notEmpty := false

	var notice []string
	if len(lines) == 0 {
		notice = p.Buffer[startingLine : endLine+1]
	} else if len(lines) > 0 {
		notice = p.Buffer[startingLine:lines[0]]
	}
	val := strings.Join(trimLeadingAndTrailingEmptyLines(notice), "\n")
	if val != "" {
		notEmpty = true
		changes.Notice = &val
	}

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
