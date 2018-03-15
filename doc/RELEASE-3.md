# ANWORK Version 3 Release Notes

## New Functionality

- The anwork\_testing repo has been deprecated in favor of an integration package in this repo.
- The anwork integration tests now contain backwards compatibility tests.
- Anwork data is now stored on disk as JSON.

## Changed Functionality

- The storage.Persister type has been changed to a generic interface. The new storage.FilePersister
  type contains the old storage.Persister functionality.
- Error messages for unknown task names have been slightly improved.

## Deprecated Functionality

There is no deprecated functionality in this release.

## Removed Functionality

There is no removed functionality in this release.
