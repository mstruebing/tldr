# Changelog

## Unreleased
### Added
* Added docker container [#50](https://github.com/mstruebing/tldr/pull/50) ([@mstruebing](https://github.com/mstruebing))
### Changed
* Only load new cache when connected to the internet and remote host is reachable [#49](https://github.com/mstruebing/tldr/pull/49) ([@mstruebing](https://github.com/mstruebing))
### Deprecated
### Removed
### Fixed
* golang ci errors
### Security
### Misc
* use go mod

## [1.1.1] - 2019-02-19
### Changed
* the chache directory now follows the XDG-standard (https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html)
* removed 1.8/1.9/1.10 from travis and only use the latest 2 versions
### Fixed
* only consider markdown files as pages (as there was an `index.json`-file added)
