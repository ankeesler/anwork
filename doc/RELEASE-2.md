# ANWORK Version 2 Release Notes

## New Functionality

- The new anwork executable is over 75% faster! Woo! This was calculated via the benchmark functions
  in the [anwork_testing](https://github.com/ankeesler/anwork_testing) repo (see BenchmarkCreate and
  BenchmarkCrud).
- The CLI commands have been significantly simplified.

## Changed Functionality

The following CLI commands have been simplified.

| Old Command | New Commands | Notes |
| --- | --- | --- |
| `task create` | `create` | See Removed Functionality below for further updates to this command. |
| `task set-*` | `set-*` | This includes both the commands that set state and priority. |
| `task delete` | `delete` | |
| `task delete-all` | `delete-all` | |
| `task show` | `show` | See Removed Functionality below for further updates to this commands. |
| `task note` | `note` | |
| `journal show-all` | `journal` | |
| `journal show` | `journal` | The `journal` command now accepts an optional task name argument. |

## Deprecated Functionality

There is no deprecated functionality in this release.

## Removed Functionality

- Long flags (of the format --flag) are no longer used. The CLI commands only accept single letter
  short flags (of the format -f).
- The `create` (formerly `task create`) command no longer takes any flags. A task's priority can be
  set via the `set-priority` command.
- The `show` (formerly `task show`) command no longer takes any flags. It always prints in the legacy
  `-s` flag format. To show a task's most recent journal entry, use the  `journal` command (e.g., "anwork journal weigh-tuna").
  To show information about a task, use the `show` command (e.g., "anwork show weigh-tuna").
- The `-n|--no-persist` flag has been removed. In order to mimic this behavior, remove the context
  file after every use of anwork.
- The `-f|--force` flag has been removed from the `reset` CLI command. In order to incite this
  behavior, simply pass a second argument of `y` to the `reset` command.
