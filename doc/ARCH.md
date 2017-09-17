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