# ANWORK

ANWORK is a personal task management system.

[![Build Status](https://travis-ci.org/ankeesler/anwork.svg?branch=master)](https://travis-ci.org/ankeesler/anwork)
[![codecov](https://codecov.io/gh/ankeesler/anwork/branch/master/graph/badge.svg)](https://codecov.io/gh/ankeesler/anwork)

Latest release: [v2](https://github.com/ankeesler/anwork/releases/tag/v2)

## Running

To get up and running with ANWORK, try the following.

0. Download Go! You can do so with Homebrew:
```
$ brew install go
```
1. Clone this repo at a location that follows the directory structure described here:
   https://golang.org/doc/code.html#Workspaces. Make sure to set your GOPATH variable accordingly!
2. Run the following command to download the dependencies needed to run anwork.
```
$ go get ./...
```
3. Run anwork with the following command.
```
$ go run cmd/anwork/main.go
```

See [CLI-OVERVIEW.md](doc/CLI-OVERVIEW.md) for full usage documentation.

## Developing

To develop on the anwork project, there are a number of dependencies that need to be installed.

0. Download Go! You can do so with Homebrew:
```
$ brew install go
```
1. The Google Protocol Buffers compiler (protoc). You can download it with Homebrew:
```
$ brew install protoc
```
2. The Go Protocol Buffer compiler plugin (protoc-gen-go). You can download it with "go get":
```
$ go get github.com/golang/protobuf/protoc-gen-go
```
Make sure that the $GOPATH/bin directory is in your $PATH.

To run the unit tests, try the following.
```
$ go get -t ./...                             # download dependencies for testing
$ go install -v github.com/onsi/ginkgo/ginkgo # install ginkgo testing binary
$ ginkgo ./...                                # run the tests with the ginkgo testing binary
```

See [anwork_testing](https://github.com/ankeesler/anwork_testing) repo for more testing.

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
