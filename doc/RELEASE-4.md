# ANWORK Version 4 Release Notes

## New Functionality

- There is now a runner package for driving the task.Manager type. This also allows for improved testing of this stuff. This package contains most of the stuff from the cmd/anwork/... packages.

## Changed Functionality

- Major update of task package.
  - The storage functionality is now baked into the task package via a task.Factory.
  - The task.Journal functionality is now baked into the task package via a task.Manager.

## Deprecated Functionality

There is no deprecated functionality in this release.

## Removed Functionality

- Protocol Buffer persistence is no longer supported. So that means v2 and v3 are no longer supported at all.
