# ANWORK

ANWORK is a personal task management system.

[![Build Status](https://travis-ci.org/ankeesler/anwork.svg?branch=feature/go-ify)](https://travis-ci.org/ankeesler/anwork)
[![codecov](https://codecov.io/gh/ankeesler/anwork/branch/feature%2Fgo-ify/graph/badge.svg)](https://codecov.io/gh/ankeesler/anwork)

Latest release: [v1](https://github.com/ankeesler/anwork/releases/tag/v1)

## Quickstart

To get up and running with ANWORK, try the following.
```
$ go run cmd/anwork/anwork.go version
```

See [CLI-OVERVIEW.md](doc/CLI-OVERVIEW.md) for full usage documentation.

## Developing

To develop on the anwork project, there are a number of dependencies that need to be installed.
1. The Google Protocol Buffers compiler (protoc). You can download it with Homebrew:
```
$ brew install protoc
```
2. The Go Protocol Buffer compiler plugin (protoc-gen-go). You can download it with "go get":
```
$ go get github.com/golang/protobuf/protoc-gen-go
```


To run the unit tests, try the following.
```
$ go test ./...
```

- See [ARCH.md](doc/ARCH.md) for architecture information.
- See [FEATURE.md](doc/FEATURE.md) for planned features.
- See [anwork_testing](https://github.com/ankeesler/anwork_testing) repo for more testing.