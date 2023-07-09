# go-changelog

[![Go Reference](https://pkg.go.dev/badge/github.com/anton-yurchenko/go-changelog.svg)](https://pkg.go.dev/github.com/anton-yurchenko/go-changelog)
[![Code Coverage](https://codecov.io/gh/anton-yurchenko/go-changelog/branch/main/graph/badge.svg)](https://codecov.io/gh/anton-yurchenko/go-changelog)
[![Go Report Card](https://goreportcard.com/badge/github.com/anton-yurchenko/go-changelog)](https://goreportcard.com/report/github.com/anton-yurchenko/go-changelog)
[![Release](https://img.shields.io/github/v/release/anton-yurchenko/go-changelog)](https://github.com/anton-yurchenko/go-changelog/releases/latest)
[![License](https://img.shields.io/github/license/anton-yurchenko/go-changelog)](LICENSE.md)

Golang package for changelog file creation/parsing

## Features

- Supports [Semantic Version](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/) Compliant
- [Common Changelog](https://common-changelog.org/) Compliant

## Manual

- Install with `go get -u github.com/anton-yurchenko/go-changelog`

### Examples

#### Create a changelog file

```golang
package main

import (
    changelog "github.com/anton-yurchenko/go-changelog"
    "github.com/spf13/afero"
)

func main() {
    c := changelog.NewChangelog()
    c.SetTitle("Changelog")
    c.SetDescription("This file contains changes of all releases")

    c.AddUnreleasedChange("fixed", []string{"Bug"})
    c.AddUnreleasedChange("added", []string{"Feature"})

    r, err := c.CreateReleaseFromUnreleasedWithURL("1.0.0", "2021-05-31","https://github.com/anton-yurchenko/go-changelog/releases/tag/v1.0.0")
    if err != nil {
        panic(err)
    }

    if err := r.AddChange("changed", "User API"); err != nil {
        panic(err)
    }
    r.AddNotice("**This release contains breaking changes**")

    if err := c.SaveToFile(afero.NewOsFs(), "./CHANGELOG.md"); err != nil {
        panic(err)
    }
}
```

#### Parse an existing changelog file

```golang
package main

import (
    "fmt"
    changelog "github.com/anton-yurchenko/go-changelog"
)

func main() {
    p, err := changelog.NewParser("./CHANGELOG.md")
    if err != nil {
        panic(err)
    }

    c, err := p.Parse()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Changelog contains %v releases", c.Releases.Len())
}
```

#### Update an existing changelog file

<details><summary>Click to expand</summary>

```golang
package main

import (
    changelog "github.com/anton-yurchenko/go-changelog"
    "github.com/spf13/afero"
)

func main() {
    p, err := changelog.NewParser("./CHANGELOG.md")
    if err != nil {
        panic(err)
    }

    c, err := p.Parse()
    if err != nil {
        panic(err)
    }

    r := c.GetRelease("1.2.1")
    if r == nil {
        panic("Release does not exists")
    }

    r.Yanked = true

    c.SaveToFile(afero.NewOsFs(), "./CHANGELOG.md")
    if err != nil {
        panic(err)
    }
}
```

</details>  

## Notes

- Releases are sorted by their [Semantic Version](https://semver.org/)
- Scopes are sorted by their importance
- `SaveToFile` will overwrite the existing file, and anything that does not match the changelog format will be omitted

## License

[MIT](LICENSE.md) © 2021-present Anton Yurchenko
