package storage

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

const (
	goodRoot = "test-data/good-root"
	badRoot  = "test-data/bad-root"

	emptyContext     = "empty-context"
	singletonContext = "singleton-context"
	badContext       = "bad-context"

	tmpRoot = "test-data/root.tmp"

	tmpEmptyContext     = "empty-context.tmp"
	tmpSingletonContext = "singleton-context.tmp"
)

func TestExists(t *testing.T) {
	data := []struct {
		name    string
		root    string
		context string
		exists  bool
	}{
		{"ContextExists", goodRoot, emptyContext, true},
		{"ContextDoesntExist", goodRoot, badContext, false},
		{"BadRoot", badRoot, emptyContext, false},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			p := Persister{datum.root}
			if p.Exists(datum.context) != datum.exists {
				t.Errorf("Root: %s, Context: %s, Expected %v but got %v",
					datum.root, datum.context, datum.exists, !datum.exists)
			}
		})
	}
}

type GoodSerializer struct {
}

func (s *GoodSerializer) Unserialize(bytes []byte) (interface{}, error) {
	return bytes, nil
}

func (s *GoodSerializer) Serialize(thing interface{}) ([]byte, error) {
	return thing.([]byte), nil
}

type BadSerializer struct {
}

func (s *BadSerializer) Unserialize(bytes []byte) (interface{}, error) {
	return nil, errors.New("Failure!")
}

func (s *BadSerializer) Serialize(thing interface{}) ([]byte, error) {
	return nil, errors.New("Failure!")
}

func validateBytes(t *testing.T, stuff interface{}, expectedBytes []byte) {
	stuffAsBytes, ok := stuff.([]byte)
	if !ok {
		t.Errorf("Expected stuff to be a []byte, but was %v", stuff)
	} else if !bytes.Equal(stuffAsBytes, expectedBytes) {
		t.Errorf("Expected %v stuff got %v", expectedBytes, stuffAsBytes)
	}
}

func TestUnpersist(t *testing.T) {
	data := []struct {
		name       string
		root       string
		context    string
		serializer Serializer
		stuff      []byte // nil for failure
	}{
		{"RootDoesntExist", badRoot, emptyContext, &GoodSerializer{}, nil},
		{"ContextDoesntExist", goodRoot, badContext, &GoodSerializer{}, nil},
		{"BadSerializer", goodRoot, emptyContext, &BadSerializer{}, nil},
		{"EmptyContext", goodRoot, emptyContext, &GoodSerializer{}, []byte{}},
		{"SingletonContext", goodRoot, singletonContext, &GoodSerializer{}, []byte{'a', 'b', 'c'}},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			t.Logf("Root: %s, Context: %s", datum.root, datum.context)

			p := Persister{datum.root}
			stuff, err := p.Unpersist(datum.context, datum.serializer)
			if datum.stuff == nil { // expecting failure...
				if err == nil {
					t.Errorf("Expected error, but got success")
				}
			} else { // expecting success...
				if err != nil {
					t.Errorf("Expected success, but got error: %s", err)
				} else {
					validateBytes(t, stuff, datum.stuff)
				}
			}
		})
	}
}

func TestPersist(t *testing.T) {
	if _, err := os.Stat(tmpRoot); !os.IsNotExist(err) {
		t.Fatalf("Error: cannot run test when tmpRoot (%s) exists", tmpRoot)
	}
	defer os.RemoveAll(tmpRoot)

	data := []struct {
		name       string
		root       string
		context    string
		serializer Serializer
		stuff      []byte // nil for failure
	}{
		{"BadSerializer", tmpRoot, tmpEmptyContext, &BadSerializer{}, nil},
		{"EmptyContext", tmpRoot, tmpEmptyContext, &GoodSerializer{}, []byte{}},
		{"SingletonContext", tmpRoot, tmpSingletonContext, &GoodSerializer{}, []byte{'a', 'b', 'c'}},
		{"SingletonContextExists", tmpRoot, tmpSingletonContext, &GoodSerializer{}, []byte{'c', 'b', 'a'}},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			t.Logf("Root: %s, Context: %s", datum.root, datum.context)

			p := Persister{datum.root}
			err := p.Persist(datum.context, datum.serializer, datum.stuff)
			if datum.stuff == nil { // expecting failure...
				if err == nil {
					t.Errorf("Expected error, but got success")
				}
			} else { // expecting success...
				if err != nil {
					t.Errorf("Expected success, but got error: %s", err)
				} else {
					stuff, err := p.Unpersist(datum.context, datum.serializer)
					if err != nil {
						t.Errorf("Failed to unpersist persisted bytes: %s", err)
					}
					validateBytes(t, stuff, datum.stuff)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	if _, err := os.Stat(tmpRoot); !os.IsNotExist(err) {
		t.Fatalf("Error: cannot run test when tmpRoot (%s) exists", tmpRoot)
	}
	defer os.RemoveAll(tmpRoot)

	data := []struct {
		name    string
		root    string
		context string
		create  bool
	}{
		{"BadRoot", tmpRoot, tmpEmptyContext, false},
		{"BadContext", tmpRoot, badContext, false},
		{"GoodContext", tmpRoot, tmpSingletonContext, true},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			t.Logf("Root: %s, Context: %s", datum.root, datum.context)

			p := Persister{datum.root}
			if datum.create {
				err := p.Persist(datum.context, &GoodSerializer{}, []byte{})
				if err != nil {
					t.Fatalf("Could not persist to context in order to delete it: %s", err)
				}
			}

			err := p.Delete(datum.context)
			if err != nil {
				t.Errorf("Got error from deleting: %s", err)
			}
		})
	}
}
