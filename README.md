# go-changelog

[![Go Reference](https://pkg.go.dev/badge/github.com/anton-yurchenko/go-changelog.svg)](https://pkg.go.dev/github.com/anton-yurchenko/go-changelog)
[![Code Coverage](https://codecov.io/gh/anton-yurchenko/go-changelog/branch/main/graph/badge.svg)](https://codecov.io/gh/anton-yurchenko/go-changelog)
[![Go Report Card](https://goreportcard.com/badge/github.com/anton-yurchenko/go-changelog)](https://goreportcard.com/report/github.com/anton-yurchenko/go-changelog)
[![Release](https://img.shields.io/github/v/release/anton-yurchenko/go-changelog)](https://github.com/anton-yurchenko/go-changelog/releases/latest)
[![License](https://img.shields.io/github/license/anton-yurchenko/go-changelog)](LICENSE.md)

Golang package for parsing a changelog file

## Overview

`go-changelog` support changelog files compliant with [Common Changelog](https://github.com/vweevers/common-changelog) and [Keep a Changelog](https://keepachangelog.com/).

## Manual

1. Install with `go get -u github.com/anton-yurchenko/go-changelog`
2. Create a parser by providing a changelog file path to it and parse the content:

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

    changelog, err := parser.Parse()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Changelog contains %v releases", len(changelog.Releases))
}
```

## License

[MIT](LICENSE.md) © 2021-present Anton Yurchenko
