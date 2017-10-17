# ANWORK FEATURE REQUESTS

This document contains feature requests for the ANWORK project. The feature requests are divided up
into sections below.

## Application
- Add note command for adding a note to a task.
- Add ability to optionally specify description and priority of task.

## CLI
- Add CLI argument stuff.
- Add support for CliArgumentType.FILE.

## Infrastructure
- Can we mandate that protocol buffer classes are named accordingly?

## Journaling

## Persistence

## Serialization
- We need a more defined way to get a class' Serializer.

## Tasks
- Write our own heap so that we can easily get the tasks in their sorted order.
- Cache the TaskManagerJournal per-task instances so we don't keep creating new objects.