# ANWORK

ANWORK is a personal task management system.

[![Build Status](https://travis-ci.org/ankeesler/anwork.svg?branch=master)](https://travis-ci.org/ankeesler/anwork)
[![codecov](https://codecov.io/gh/ankeesler/anwork/branch/master/graph/badge.svg)](https://codecov.io/gh/ankeesler/anwork)

Latest release: [v1](https://github.com/ankeesler/anwork/releases/tag/v1)

## Quickstart

To get up and running with ANWORK, try the following.
```
$ ./gradlew install
$ PATH="$PATH:$PWD/build/install/anwork/bin"
$ anwork task create 'hello world'
$ anwork task show
```

See [CLI-OVERVIEW.md](doc/CLI-OVERVIEW.md) for full usage documentation.

## Developing

To run the unit tests, try the following.
```
$ ./gradlew check
```

To run the smoke tests, try the following.
```
$ ./gradlew smoke
```

- See [ARCH.md](doc/ARCH.md) for architecture information.
- See [FEATURE.md](doc/FEATURE.md) for planned features.
- See [anwork_testing](https://github.com/ankeesler/anwork_testing) repo for more testing.