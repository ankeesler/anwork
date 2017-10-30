# ANWORK ARCHITECTURE

This document contains the software architecture for the ANWORK project. The information below
contains not only core software design explanations, but also infrastructure notes and best
practices.

## General Design Principles

There are a lot of places where we use the java.lang.String type in a weak-typed manner. This is an
experiment. We are curious whether or not using a java.lang.String is a better way forward than
using a java.lang.Object or generic type.

## Core

The com.marshmallow.anwork.core package contains utilities and interfaces that apply to the whole
ANWORK project.

### Serialization

The com.marshmallow.anwork.core.Serializ* interfaces are a framework that can take an object and
turn it into an array of bytes. I did not use the java.io.Serializable framework here because that
would be too easy. There are 2 core ideas to understand here.
- A **Serializable** is something that can be turned into bytes (by a Serializer) via some medium.
  This medium may be XML, Google Protocol Buffers, a java.lang.String, a custom implementation, etc.
- A **Serializer** is something that does the work of taking an object and turning it into bytes.

There is currently one com.marshmallow.anwork.core.Serializer implementation. This implementation
(com.marshmallow.anwork.core.ProtobufSerializer) serializes objects via Google Protocol Buffers.

The com.marshmallow.anwork.core.Factory interface is a utility interface for creating default
instances of objects to be used in unserialization.

It is **highly** advised that com.marshmallow.anwork.core.Serializable objects offer a static method
or field by which clients can access a com.marshmallow.anwork.core.Serializer for that object. In
the future perhaps we should create a factory for this purpose.

### Persistence

A com.marshmallow.anwork.core.Persister is something that uses a
com.marshmallow.anwork.core.Serializer to persist something to a context. A context is a purposely
vague piece of data that is implementation specific. It may be a location on disk, or a URI, or a
database key, etc. See General Design Principles for more information about vague pieces of data
such as this.

There is currently one com.marshmallow.anwork.core.Persister implementation. This implementation
(com.marshmallow.anwork.core.FilePersister) persists serialized objects to disk via a java.io.File.

## Journal

The journal framework provides a way to maintain a history of events. The framework goes further to
support the use case of getting a specific collection of events with a common quality in the order
that they occured. Here are the major interfaces.
- A **JournalEvent** is an object that describes some activity that has happened.
- A **Journal** is an object that maintains an ordering of events (JournalEvent's) that have
  happened.
- A **Journaled** is an object that has an associated journal.
- A **MultiJournaled** is a Journaled that has multiple Journal instances that can be acquired based
  on a key. This key type is implementation specific in order to support all types as a key.

## CLI

### Public Interface

The CLI framework is a (large) add-on to the application layer. It is a general framework for
interacting with a Java application at the command line. Here are the main public concepts. These
classes are located in the com.marshmallow.anwork.app.cli package.
- A *Cli* is the entry point object for creating a CLI for an application. A client can create an
  instance of this object and create their CLI API via the root list (See *Cli#getRoot*). The
  method *Cli#parse* is what actually does the parsing of command line arguments.
- A *Command* is a keyword passed on the command line that runs some *Action*. A *Command* may be
  something like "commit" in "git commit." An *Action* is just some Java code that runs when this
  *Command* is passed. An *ActionCreator* is a simple interface that creates *Action*s.
- A *List* is a collection of CLI *Command*'s. *List*s allow us to organize *Command*s into
  different sub APIs. For example, "git submodule" is a *List* because it contains a number of
  *Command*s (update, add, remove, etc.).
- A *Flag* is an option passed to a *Command* or *List*. It is of the format -f (referred to as a
  "short flag") or --file (referred to as a "long flag"). An example of a *Flag* is "-f" in
  "make -f file.mak".
- An *Argument* is a thing that is passed to a *Command* or a *Flag*. *Argument*s have a backing
  Java type that allows clients to easily convert from command line to Java runtime. An example of
  an *Argument* would be "file.mak" in "make -f file.mak" or "25" in "tail -n 25".
- A *Visitor* is an object that can iterate through the CLI data structure and be notified of
  *Flag*s, *Command*s, and *List*s. See *Cli#visit* for the main use of this class.

Most of the main concepts have sub-interfaces that start with "Mutable". These sub-interfaces allow
for editing of the main concepts, i.e., they offer setters whereas the super-interfaces only offer
getters.

Note that the same "short flag" cannot be added to two lists or a list and a command in the same
CLI hiearchy. For example, consider the command line invocation "git -a whatever log -a foo". Upon
the calling of the *Action* for the "log" *Command*, there would be ambiguity about what value the
"-a" *Flag* holds.

A CLI API for an application can also be specified via an XML schema. A *CliXmlReader* is the class
that reads in an XML stream and returns a *Cli* object.

Given a *Cli* instance, one can use the *DocumentationGenerator* framework to generate
documentation of different formats. The singleton *DocumentationGeneratorFactory* will generate
*DocumentationGenerator*s based on a *DocumentationType* (like Github markdown, text, etc.).  

### Internal Implementation

Most of this package includes implementation classes ending in "Impl". For example, the
implementation of *MutableArgument* is named *ArgumentImpl*.

This CLI implementation uses a tree to store *List*s (see *ListImpl*) and *Command*s (see
*CommandImpl*). When *Cli#parse* is called, the *ListImpl* class lazily initializes a
*ParseContext* instance in order to keep track of what *Flag*s, *List*s, and *Command*s are
available on the current *ListImpl*. The *Cli#parse* method works down the tree,
collecting arguments and flags as it goes. If an unknown argument to a list is encountered, then
the parsing stops and an IllegalArgumentException is throw. If a valid *Command* is encountered,
the arguments are validated (see *CommandImpl#validateContext) and the *Action* for that *Command*
is run (see *CommandImpl#runActionFromContext).

## Test

All tests must go in the src/test source set. Here are naming conventions for tests.
- Tests for classes should live in the analagous package name in the src/test source set. A
  package's analagous test package is simple the package's name plus *.test*. For example, the
  package com.marshmallow.anwork.tuna would have a test package com.marshmallow.anwork.tuna.test.
- Test classes should be named starting with the name of the class being tested and ending with
  "Test." For example, the class TunaFish would have a test class called TunaFishMarlin.
- Superclasses of tests should be begin with *Base* so that the gradle test action can filter these
  tests out. For example, a base test class for an interface Marlin should be called BaseMarlinTest.
  See build.gradle for more information.
- Junit suite definitions should begin with *All* and end in *Tests*. Each test package should have
  a suite that begins with the name of the last segment of the package name. For example, in the
  package com.marshmallow.anwork.foo.test, there should be a Junit suite class named AllFooTests.
- The AllTests class should be kept up to date to run all tests for the ANWORK project.