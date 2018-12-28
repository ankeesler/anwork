# ANWORK Version 8 Release Notes

We are experimenting with using Pivotal Tracker to keep track of work! Check it out [here](https://www.pivotaltracker.com/n/projects/2230869)!

## New Functionality
- Command shortcuts like "sr" for "set-running" or "n" for "note" or "c" for "create".
- API supports authentication! Each API deployment is accessible via a single RSA key/secret.
- API documentation is generated. And added to release script.
- There is an integration test for the API+client.

## Changed Functionality
- ANWORK architecture now uses repository abstraction to represent a collection of `task.Task`'s and `task.Event`'s.
- ANWORK architecture now separated repository layer (`task.Repo`) with manager (`manager.Manager`) layer.
- Bug fixed where you couldn't set tasks to "Ready" with the API.
- The `api` package was refactored to be more correctly tested. Now it only provides an `http.Handler` that encapsulates the ANWORK API.

## Deprecated Functionality

## Removed Functionality
