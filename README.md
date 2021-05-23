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
- [Common Changelog](https://github.com/vweevers/common-changelog) Compliant

## Manual

- Install with `go get -u github.com/anton-yurchenko/go-changelog`


### Examples

#### Create a changelog file

```golang
package main

import (
    changelog "github.com/anton-yurchenko/go-changelog"
)

func main() {
    changelogContent, err := changelog.NewChangelog("Changelog", "")
    if err != nil {
        panic(err)
    }

    changelogContent.AddUnreleasedChanges("Fixed", []string{"Bug"})

    changelogContent.AddRelease("1.0.0", "https://github.com/anton-yurchenko/go-changelog/releases/tag/v1.0.0", "2021-05-19")

    changelogContent.AddReleaseChanges("1.0.0", "Added", []string{
        "Feature A",
        "Feature B",
    })

    changelogContent.SaveToFile("./CHANGELOG-FORMATTED.md")
    if err != nil {
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
    parser, err := changelog.NewParser("./CHANGELOG.md")
    if err != nil {
        panic(err)
    }

    changelogContent, err := parser.Parse()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Changelog contains %v releases", len(changelogContent.Releases))
}
```

#### Update an existing changelog file

<details><summary>Click to expand</summary>

```golang
package main

import (
    changelog "github.com/anton-yurchenko/go-changelog"
)

func main() {
    parser, err := changelog.NewParser("./CHANGELOG.md")
    if err != nil {
        panic(err)
    }

    changelogContent, err := parser.Parse()
    if err != nil {
        panic(err)
    }

    if err := changelogContent.YankRelease("1.2.1"); err != nil {
        panic(err)
    }

    changelogContent.SaveToFile("./CHANGELOG.md")
    if err != nil {
        panic(err)
    }
}
```

</details>  

## Notes

- Releases are sorted by their [Semantic Version](https://semver.org/)

## License

[MIT](LICENSE.md) © 2021-present Anton Yurchenko
