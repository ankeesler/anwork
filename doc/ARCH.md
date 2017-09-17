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

There are some simple objects provided as base classes for these types. These objects are trivial to
understand. See **BaseJournalEvent** and **BaseJournal**.

The most complicated type in this package is the **FilteredJournal**. It is a Journal that can
make itself yield different events depending on a filter applied to another Journal.

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