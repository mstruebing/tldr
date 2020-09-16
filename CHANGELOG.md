# Changelog

## Unreleased

### Added

### Changed

### Deprecated

### Removed

### Fixed

### Security

### Misc

## [1.2.2] - 2020-09-16

### Fixed

-   Get port from scheme [#55](https://github.com/mstruebing/tldr/pull/55) ([@mstruebing](https://github.com/mstruebing))

## [1.2.1] - 2020-07-29

-   Only check remote is reachable when ttl timeout. [#52](https://github.com/mstruebing/tldr/pull/52) ([@wudong](https://github.com/wudong))

## [1.2.0] - 2020-03-06

### Added

-   Added docker container [#50](https://github.com/mstruebing/tldr/pull/50) ([@mstruebing](https://github.com/mstruebing))

### Changed

-   Only load new cache when connected to the internet and remote host is reachable [#49](https://github.com/mstruebing/tldr/pull/49) ([@mstruebing](https://github.com/mstruebing))
-   Switched to go mod [47](https://github.com/mstruebing/tldr/pull/47) ([@mstruebing](https://github.com/mstruebing))

### Fixed

-   golang ci errors

## [1.1.1] - 2019-02-19

### Changed

-   the chache directory now follows the XDG-standard (https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html)
-   removed 1.8/1.9/1.10 from travis and only use the latest 2 versions

### Fixed

-   only consider markdown files as pages (as there was an `index.json`-file added)
