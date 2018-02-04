// This package contains functionality for storing anwork data to file. This package is comprised of
// two types that allow data to be persisted between anwork uses.
//
// A Serializable is an interface that represents a type that marshal and unmarshal itself to and
// from an array of bytes, respectively.
//
// A Persister is a type that is able to take Serializable objects and read or write them to some
// persistent store. A Persister cares about a special concept called a "context." A context is
// a way to specify where the data is to be stored. This allows for users of anwork to have multiple
// different data stores depending on what they are currently working on (i.e., an "at home" to-do
// list versus an "at work" to-do list).
package storage

// An instance of this object can write any object to a persistent store (via Persist) and read it
// back from the store (via Unpersist). It uses the idea of a "context" to specify different areas
// of the persistent store. A "context" can exist or not (see Exists) and can also be deleted (via
// Delete).
type Persister interface {
	// Store the provided serializable object in the context.
	Persist(context string, serializable Serializable) error
	// Read the provided serializable object from the context.
	Unpersist(context string, serializable Serializable) error
	// Return whether or not the specific context has been created. It is expected that if Exists
	// returns true, then Unpersist will return some sort of (valid or invalid) object.
	Exists(context string) bool
	// Delete the context such that Exists will return false.
	Delete(context string) error
}
