# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.3] - 2019-06-28
### Added
- script release.sh to make automatic releases
### Changed 
- All updated entries are saved as 'billable'
- Allow to ask for `-version` of script without providing the `-api-key`

## [0.1.2] - 2019-06-22
### Added 
- Option to not display colorful terminal output.
- more readable details of what is happening

## [0.1.1] - 2019-06-21
### Added 
- Optimization: no more updates when nothing changes in the entry

## [0.1.0] - 2019-06-21
### Added 
- different strategies for rounding last entry remaining time with new `-remainig` argument
- corrected way all arguments works

## [0.0.4] - 2019-06-20
### Added
- correct handling of overlapping entries
- written tests for basic behavior 

## [0.0.3] - 2019-06-20
### Added
- `-debug` to print out detailed logs of go-toggl requests. By default it's more silent.
- this CHANGELOG.md file to repository 

## [0.0.2] - 2019-06-19
### Added
- `-api-key` argument option to provide the key instead of environment variable 

## [0.0.1] - 2019-06-19
### Added
- first working version