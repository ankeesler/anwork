package storage

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
)

// This object does the persisting of data to some file. An instance of this object can store any
// object to a file that is able to be "serialized." An object is able to be "serialized" if it
// implements the Serializable interface.
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

// Save a Serializable object to a context, or return a non-nil error iff there is a problem.
func (p *Persister) Persist(context string, serializable Serializable) error {
	if !p.Exists(context) {
		p.create(context)
	}

	bytes, err := serializable.Serialize()
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

// Load a Serializable object from a context, or return a non-nil error iff there is a problem.
func (p *Persister) Unpersist(context string, serializable Serializable) error {
	if !p.Exists(context) {
		return errors.New("Context does not exist: " + context)
	}

	path := path.Join(p.Root, context)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return serializable.Unserialize(bytes)
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
