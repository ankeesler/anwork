# ANWORK

ANWORK is a personal task management system.

[![Build Status](https://travis-ci.org/ankeesler/anwork.svg?branch=master)](https://travis-ci.org/ankeesler/anwork)

## Quickstart

To get up and running with ANWORK, try the following.
```
$ ./gradlew install
$ PATH="$PATH:$PWD/build/install/anwork/bin"
$ anwork task create 'hello world'
$ anwork task show
```

See [CLI.md](doc/CLI.md) for full CLI documentation.

## Developing

To run the stage 1 tests, try the following.
```
$ ./gradlew check
```

To run the stage 1 and 2 tests, try the following.
```
$ ./gradlew smoke
```

- See [ARCH.md](doc/ARCH.md) for architecture information.
- See [FEATURE.md](doc/FEATURE.md) for planned features.
