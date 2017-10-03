# ANWORK FEATURE REQUESTS

This document contains feature requests for the ANWORK project. The feature requests are divided up
into sections below.

## Application
- Add note command for adding a note to a task.
- Create Anwork's CLI XML document.

## CLI
- Add CLI documentation generator.
- Add CLI argument stuff.
- The CLI usage printing is ugly. Write a real test for this.
- Create a framework where a CLI can be supplied by an XML file and turned into
  source creation of CLI, tests, and documentation.
- Add a CliActionCreator interface.
- Remove the usage functionality from the CliNodeImpl class and use a CliVisitor instead.

## Infrastructure
- Can we mandate that protocol buffer classes are named accordingly?
- We need to figure out what the installed "thing" is above and make sure it has the right classes with it.
- Add in documentation generation in build process - regeneration of CLI.md needs to be added to documentation
  gradle tasks and we need for build to depend on it. Documentation before code!

## Journaling

## Persistence

## Serialization
- We need a more defined way to get a class' Serializer.