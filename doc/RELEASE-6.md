# ANWORK Version 6 Release Notes

## New Functionality
- ANWORK version includes the git hash and date at which it was built.
- ANWORK release process has a script to help automate it.
  - runner.Version is bumped.
  - git tag is added.
  - Backwards compat context is added in integration/data/.
  - CLI documentation is regenerated.
  - README is updated with latest release information.
- By default, anwork places contexts in $HOME/.anwork/.

## Changed Functionality
- The `api` package was refactored to be more correctly tested. A `service`
  package was added to encapsulate the anwork API HTTP service functionality.
- The _Waiting_ task state has been renamed to _Ready_. This is so that the
  task states look more like classical kernel scheduler states.

## Deprecated Functionality

## Removed Functionality

