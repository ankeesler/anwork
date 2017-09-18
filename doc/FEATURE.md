# ANWORK FEATURE REQUESTS

This document contains feature requests for the ANWORK project. The feature requests are divided up
into sections below.

## Application
- Add note command for adding a note to a task.
- Implement "anwork journal" CLI commands.

## CLI
- Add CLI documentation generator.
- Add CLI argument stuff.
- The CLI usage printing is ugly. Write a real test for this.
- Create a framework where a CLI can be supplied by an XML file and turned into
  source creation of CLI, tests, and documentation.

## Infrastructure
- Can we mandate that protocol buffer classes are named accordingly?
- Add some sort of gradle dependency on protoc so that users don't have to download this binary?

## Journaling
- Implement the Journaled methods on TaskManager.
- Add serialization to the journal object model.

## Persistence