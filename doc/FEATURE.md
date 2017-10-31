# ANWORK FEATURE REQUESTS

This document contains feature requests for the ANWORK project. The feature requests are divided up
into sections below.

## Application

## CLI
- Add support for CliArgumentType.FILE.
- Add examples into the CLI schema (an example should have "text" and a "description" of what the
  example actually does).
- The CLI argument outputs look terrible. Let's fix these.

## Infrastructure
- Can we mandate that protocol buffer classes are named accordingly?

## Journaling

## Persistence

## Serialization
- We need a more defined way to get a class' Serializer.

## Tasks
- Write our own heap so that we can easily get the tasks in their sorted order.
- Cache the TaskManagerJournal per-task instances so we don't keep creating new objects.