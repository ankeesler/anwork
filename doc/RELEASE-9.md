# ANWORK Version 9 Release Notes

We are experimenting with using Pivotal Tracker to keep track of work! Check it out [here](https://www.pivotaltracker.com/n/projects/2230869)!

## New Functionality
- A SQL `task.Repo` is now available for use. No more losing data when a container gets rerolled! Woo hoo!
- The ANWORK service has been updated to pull in a Cloud Foundry database service if running in that environment.
- API profiling is available on endpoint `/debug/pprof`.

## Changed Functionality
- The functions `task.Repo.DeleteTask` and `task.Repo.DeleteEvent` were updated to return no error on unknown entity.
- The use of a random `io.Reader` was updated in the `api/auth` code so it should fall over less.

## Deprecated Functionality

## Removed Functionality
