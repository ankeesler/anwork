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

## CLI

Read more about the CLI [here](doc/CLI-OVERVIEW.md).

See [CLI.md](doc/CLI.md) for a full overview of the CLI.

Here is an example ANWORK CLI session.
```bash
$ anwork create take-out-trash
$ anwork create buy-groceries-for-dinner
$ anwork show
RUNNING tasks:
BLOCKED tasks:
READY tasks:
  take-out-trash (8)
  buy-groceries-for-dinner (9)
FINISHED tasks:
$ anwork set-blocked take-out-trash
$ anwork set-running buy-groceries-for-dinner
$ anwork show
RUNNING tasks:
  buy-groceries-for-dinner (9)
BLOCKED tasks:
  take-out-trash (8)
READY tasks:
FINISHED tasks:
$ anwork set-finished buy-groceries-for-dinner
$ anwork show
...
```

## API

See [API.md](doc/API.md) for full overview of the API.
