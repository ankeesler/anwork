# ANWORK

ANWORK is a personal task management system.

[![Build Status](https://travis-ci.org/ankeesler/anwork.svg?branch=master)](https://travis-ci.org/ankeesler/anwork)
[![codecov](https://codecov.io/gh/ankeesler/anwork/branch/master/graph/badge.svg)](https://codecov.io/gh/ankeesler/anwork)

Latest release: [v4](https://github.com/ankeesler/anwork/releases/tag/v4)

## Running

To get up and running with ANWORK, run the bootstrap script!

1. Set your GOPATH environmental variable accordingly.
```
$ export GOPATH=...
```
2. Run the bootstrap.sh script at the root of this repo. If you have `curl` installed already, you can do this.
```
$ curl https://raw.githubusercontent.com/ankeesler/anwork/master/bootstrap.sh | bash
```
3. Run `anwork`.
```
$ $GOPATH/bin/anwork
$ ...
```

See [CLI-OVERVIEW.md](doc/CLI-OVERVIEW.md) for full usage documentation.

## Developing

To develop on the anwork project, run the bootstrap.sh script as above, and then you should be good to go!

To run all of the tests, try the following (assuming you are in the root of the repo and `ginkgo` is on your $PATH).
```
$ ginkgo -r .
```

To run the tests in a particular package, try the following.
```
$ ginkgo path/to/package
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
