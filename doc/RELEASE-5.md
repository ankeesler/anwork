# ANWORK Version 5 Release Notes

## New Functionality
- Beta support for ANWORK service.
  - See go doc in `api` package for details on the API.
  - Run anwork CLI with `ANWORK_API_ADDRESS` env var set to service address in order to communicate to service with CLI.
  - Run integration test with `ANWORK_TEST_RUN_WITH_API` env var set in order to test against API.
  - There is a known bug regarding setting tasks to the Waiting state.
- Updated the repo to contain a [bootstrap.sh](../ci/bootstrap.sh) script to get started with using/developing anwork. See [README.md](../README.md) for more details.

## Changed Functionality
- Move `task.ManagerFactory.Reset()` to `task.Manager.Reset()`.
- Some of the `task.Manager` APIs now return errors (e.g., `task.Manager.Delete()`).
- Never return nil from `task.Manager.Tasks()`.
- Print out the CLI commands in the integration test when run verbosely.
- The method `task.Manager.DeleteEvent()` was added.

## Deprecated Functionality

## Removed Functionality
