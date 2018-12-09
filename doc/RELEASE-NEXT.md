# ANWORK Version NEXT Release Notes

We are experimenting with using Pivotal Tracker to keep track of work! Check it out [here](https://www.pivotaltracker.com/n/projects/2230869)!

## New Functionality
- GA ANWORK service support.
  - There is an integration test for the API+client.
  - API documentation is generated. And added to release script.
  - Basic auth is used. Or more than basic auth, idk.
- The API respects the context flag.
- Use a SQL store for the backing datastore.
- Ability to add note on top of state changes.
- Show last note in the "show" view.
- Scheduling something with a deadline and have it automatically prioritized would be really nice.
- Instead of '@' for a task ID prefix, use '.'.

## Changed Functionality
- You can't set tasks to "Ready" on the API...
- The `api` package was refactored to be more correctly tested. A `service`
  package was added to encapsulate the anwork API HTTP service functionality.

## Deprecated Functionality

## Removed Functionality
