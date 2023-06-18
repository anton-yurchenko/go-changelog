# Changelog

## [1.0.5] - 2023-07-18

### Changed

- Upgrade to GoLang 1.20
- Update dependencies

## [1.0.4] - 2021-10-18

### Fixed

- Empty scopes won't be included during parsing

## [1.0.3] - 2021-09-09

### Changed

- Changelog scopes are sorted by priority (from a point of view of a developer)

## [1.0.2] - 2021-08-22

### Fixed

- Wrong package name caching by importers

## [1.0.1] - 2021-08-22 [YANKED]

### Changed

- Update to GoLang 1.17
- Update depenedencies

## [1.0.0] - 2021-05-31

### Added

- Output to string methods on changelog struct and its components
- Support YANKED releases
- Save changelog to file method
- Create new changelog
- Add/Change changelog components like releases/title/unreleased/...

## [0.0.2] - 2021-05-22

### Added

- Parse release dates

### Changed

- Enforce a date regex, allow matching up to year 3000

### Fixed

- Multiple parser bugs

## [0.0.1] - 2021-05-19

_Initial release_

[1.0.5]: https://github.com/anton-yurchenko/go-changelog/compare/v1.0.4...v1.0.5
[1.0.4]: https://github.com/anton-yurchenko/go-changelog/compare/v1.0.3...v1.0.4
[1.0.3]: https://github.com/anton-yurchenko/go-changelog/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/anton-yurchenko/go-changelog/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/anton-yurchenko/go-changelog/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/anton-yurchenko/go-changelog/compare/v0.0.2...v1.0.0
[0.0.2]: https://github.com/anton-yurchenko/go-changelog/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/anton-yurchenko/go-changelog/releases/tag/v0.0.1
