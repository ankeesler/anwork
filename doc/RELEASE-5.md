# ANWORK Version NEXT Release Notes

## New Functionality
- Beta support for ANWORK service.
  - See go doc in api package for details on the API.
  - Missing many features, contains many bugs, but you can create and show tasks with the service!
  - Run anwork CLI with `ANWORK\_API\_ADDRESS` env var set to service address in order to communicate to service with CLI.
  - Run integration test with `ANWORK\_TEST\_RUN\_WITH\_API` env var set in order to test against API.
- By default, anwork places contexts in $HOME/.anwork/.

## Changed Functionality
- Move task.ManagerFactory.Reset() to task.Manager.Reset().
- Some of the task.Manager APIs now return errors (e.g., `task.Manager.Delete()`).
- Never return nil from `task.Manager.Tasks()`.
- Print out the CLI commands in the integration test when run verbosely.

## Deprecated Functionality

## Removed Functionality