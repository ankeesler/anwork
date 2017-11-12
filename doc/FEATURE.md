# ANWORK FEATURE REQUESTS

This document contains feature requests for the ANWORK project. The feature requests are divided up
into sections below.

## Application
- Add the ability for a Task to have a category and display that category on the screen.
- Add the ability to specify variables like "finished" and "running" in "task-name" CLI arguments.
  Also, add the ability to specify task ids and officially test this. It can be like
    "set-running" $finished
  which would set all of the tasks that are finished as running.
- We need to add a version number to every package! And a version command.
- By default, the persistence context should be the home directory.

## CLI
- Add support for CliArgumentType.FILE.

## Infrastructure
- Can we mandate that protocol buffer classes are named accordingly?
- What is this mysterious failure on Travis with the Smoketest?

## Journaling

## Persistence

## Serialization
- We need a more defined way to get a class' Serializer.

## Tasks