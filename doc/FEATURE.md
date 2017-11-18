# ANWORK FEATURE REQUESTS

This document contains feature requests for the ANWORK project. The feature requests are divided up
into sections below.

## Application
- Add the ability for a Task to have a category and display that category on the screen.
- Add the ability to specify variables like "finished" and "running" in "task-name" CLI arguments.
  Also, add the ability to specify task ids and officially test this. It can be like
    "set-running" $finished
  which would set all of the tasks that are finished as running.
- We need to add a version number to every package!
- By default, the persistence context should be the home directory.
- Show the summary of a task. For example, the amount of times we moved from running to blocked and
  running to waiting. Basically the amount of times we switched out of running to something other
  than finished. We could create some sort of "stat" index.
- We should really reorganize the CLI API. Do we really need the "journal" and "task" sub-lists?
- We should use the new java.time package! It is much better and can fix our "elapsed date" issue.

## CLI
- Add support for CliArgumentType.FILE.

## Infrastructure
- Can we mandate that protocol buffer classes are named accordingly?

## Journaling

## Persistence

## Serialization
- We need a more defined way to get a class' Serializer.

## Tasks