# ANWORK Version NEXT Release Notes

## New Functionality
- Beta support for ANWORK service.
  - See go doc in api package for details on the API.
  - Run anwork CLI with `ANWORK\_API\_ADDRESS` env var set to service address in order to communicate to service with CLI.
  - Run integration test with `ANWORK\_TEST\_RUN\_WITH\_API` env var set in order to test against API.
- Updated the repo to contain a [bootstrap.sh](../ci/bootstrap.sh) script to get started with using/developing anwork. See [README.md](../README.md) for more details.

## Changed Functionality
- Move task.ManagerFactory.Reset() to task.Manager.Reset().
- Some of the task.Manager APIs now return errors (e.g., `task.Manager.Delete()`).
- Never return nil from `task.Manager.Tasks()`.
- Print out the CLI commands in the integration test when run verbosely.

## Deprecated Functionality

## Removed Functionality
