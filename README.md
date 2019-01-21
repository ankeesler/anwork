# ANWORK

ANWORK is a multitasking management system.

[![Build Status](https://travis-ci.org/ankeesler/anwork.svg?branch=master)](https://travis-ci.org/ankeesler/anwork)
[![codecov](https://codecov.io/gh/ankeesler/anwork/branch/master/graph/badge.svg)](https://codecov.io/gh/ankeesler/anwork)

Latest release: [v9](https://github.com/ankeesler/anwork/releases/v9)

Upcoming features/bugs/refactors: [Tracker Project](https://www.pivotaltracker.com/n/projects/2230869)

## Philosophy

- As a human, multitasking is hard.
- As a computer, multitasking is easy.
- If humans apply some concepts from multitasking operating systems, we may be able to improve our multitasking abilities.

## Running

To get up and running with ANWORK...
- If you are running on darwin, download the latest binary [here](https://github.com/ankeesler/anwork/releases/latest)
- If you are running on linux, I am sure you know how to build Go code, so build the `./cmd/anwork` package
- If you are running on windows, ..., meh

See [CLI-OVERVIEW.md](doc/CLI-OVERVIEW.md) for full usage documentation.

## Developing

To develop on the anwork project...
- Make sure you have at least go version 1.11 - I recommend using [`brew`](https://brew.sh/) (see command below).
  - `$ brew install go`
- To run all the tests...
  - `$ ./ci/test.sh`
- To run a single package of tests...
  - `$ ginkgo ./path/to/package`

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
the `go doc` command.

| Directory | Use |
| --- | --- |
| ci | Scripts used in the Travis Continuous Integration jobs |
| doc | Documentation about the anwork executable, e.g., a quick start guide |
| integration | Integration tests for the anwork executable/API |
