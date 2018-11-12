# ANWORK Version NEXT Release Notes

## New Functionality
- GA ANWORK service support.
  - There is an integration test for the API+client.
  - API documentation is generated. And added to release script.
- ANWORK "archive" command deletes finished tasks.
- The API respects the context flag.
- Use a SQL store for the backing datastore.
- Post anwork binaries on github release page along with release notes and update README to point there.
- Ability to add note on top of state changes.
- Show last note in the "show" view.
- A web view that listens to the JSON data and updates whenever it changes!
- Ability to rename task.
- Scheduling something with a deadline and have it automatically prioritized would be really nice.

## Changed Functionality
- You can't set tasks to "Ready" on the API...
- The `api` package was refactored to be more correctly tested. A `service`
  package was added to encapsulate the anwork API HTTP service functionality.

## Deprecated Functionality

## Removed Functionality
