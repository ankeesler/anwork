# ANWORK

ANWORK is a personal task management system.

## TODO

### Application
- Add note command for adding a note to a task.

### Persistence
- The serialization system should support escaping certain characters.
- The #serializer methods should be #getSerializer so that the methods start with
  a verb.
- There should probably be a #clear method on Persister to erase a context.
- We should use Protocol Buffers for serialization of data. We just should.

### Infrastructure

### CLI
- Add CLI documentation generator.
- Add CLI argument stuff.
- The CLI usage printing is ugly. Write a real test for this.
- Create a framework where a CLI can be supplied by an XML file and turned into
  source creation of CLI, tests, and documentation.
