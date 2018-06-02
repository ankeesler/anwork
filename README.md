# ANWORK

ANWORK is a personal task management system.

[![Build Status](https://travis-ci.org/ankeesler/anwork.svg?branch=master)](https://travis-ci.org/ankeesler/anwork)
[![codecov](https://codecov.io/gh/ankeesler/anwork/branch/master/graph/badge.svg)](https://codecov.io/gh/ankeesler/anwork)

Latest release: [v4](https://github.com/ankeesler/anwork/releases/tag/v4)

## Running

To get up and running with ANWORK, try the following.

0. Download Go! You can do so with Homebrew:
```
$ brew install go
```
1. `go get` this repo.
```
$ go get github.com/ankeesler/anwork/cmd/anwork
```
2. Run `anwork`.
```
$ $GOPATH/bin/anwork
$ ...
```

See [CLI-OVERVIEW.md](doc/CLI-OVERVIEW.md) for full usage documentation.

## Developing

To develop on the anwork project, there is only one dependency that needs to be installed.

0. Download Go! You can do so with Homebrew:
```
$ brew install go
```

To run the unit tests, try the following.
```
$ go get -t ./...                          # download dependencies for testing
$ go install github.com/onsi/ginkgo/ginkgo # install ginkgo testing binary
$ ginkgo -r .                              # run the tests with the ginkgo testing binary
```

### Style

Thanks to the use of Go, there are not a lot of style conventions to note here. Here are the few
formatting rules used in this codebase.
1. Use `go fmt` for formatting all code.
2. No line should extend past column 100.
3. If function arguments would go past the 100 character limit, they should be wrapped to the next
   line.

### Directory Structure

Again, the directory structure is pretty simple. Every directory is a Go source directory except for
the following that are documented here. Information can be found on the Go source directories using
the "go doc" command.

| Directory | Use |
| --- | --- |
| ci | Scripts used in the Travis Continuous Integration jobs |
| doc | Documentation about the anwork executable, e.g., a quick start guide |
| integration | Integration tests for the anwork executable |
