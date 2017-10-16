# ANWORK FEATURE REQUESTS

This document contains feature requests for the ANWORK project. The feature requests are divided up
into sections below.

## Application
- Add note command for adding a note to a task.

## CLI
- Add CLI documentation generator.
- Add CLI argument stuff.
- The CLI usage printing is ugly. Write a real test for this.
- Remove the usage functionality from the CliNodeImpl class and use a CliVisitor instead.
- Add support for CliArgumentType.FILE.
- Update CLI ARCH.md text.
- Add test to make sure changing the name/short flag of a CLI node/flag doesn't mess anything up.

## Infrastructure
- Can we mandate that protocol buffer classes are named accordingly?

## Journaling

## Persistence

## Serialization
- We need a more defined way to get a class' Serializer.

## Tasks
- Write our own heap so that we can easily get the tasks in their sorted order.
- Cache the TaskManagerJournal per-task instances so we don't keep creating new objects.