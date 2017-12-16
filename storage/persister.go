package storage

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
)

// This object does the persisting of data to some file directory. An instance of this object can
// store any object to a file that is able to be "serialized." An object is able to be "serialized"
// if it has an associated Serializer. Note that each context holds a single serializable object.
type Persister struct {
	Root string
}

func (p *Persister) create(context string) error {
	err := os.MkdirAll(p.Root, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}

	path := path.Join(p.Root, context)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	f.Close()

	return nil
}

// Returns whether or not this context exists, i.e., whether or not there is an existing file for
// this context.
func (p *Persister) Exists(context string) bool {
	path := path.Join(p.Root, context)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Load an object from a context, or return a non-nil error iff there is a problem.
func (p *Persister) Unpersist(context string, serializer Serializer) (interface{}, error) {
	if !p.Exists(context) {
		return nil, errors.New("Context does not exist: " + context)
	}

	path := path.Join(p.Root, context)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return serializer.Unserialize(bytes)
}

// Save an object to a context, or return a non-nil error iff there is a problem. If the context
// does not exist, it will be created.
func (p *Persister) Persist(context string, serializer Serializer, thing interface{}) error {
	if !p.Exists(context) {
		p.create(context)
	}

	bytes, err := serializer.Serialize(thing)
	if err != nil {
		return err
	}

	path := path.Join(p.Root, context)
	err = ioutil.WriteFile(path, bytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Completely delete a context, wiping out all of its data. If the context doesn't exist, this
// function will exit silently.
func (p *Persister) Delete(context string) error {
	var err error = nil
	if p.Exists(context) {
		path := path.Join(p.Root, context)
		err = os.Remove(path)
	}
	return err
}
