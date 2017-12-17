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

type GoodSerializable struct {
	ExpectedBytes []byte
	actualBytes   []byte
}

func (s *GoodSerializable) Serialize() ([]byte, error) {
	return s.ExpectedBytes, nil
}

func (s *GoodSerializable) Unserialize(bytes []byte) error {
	s.actualBytes = bytes
	return nil
}

type BadSerializable struct {
}

func (s *BadSerializable) Serialize() ([]byte, error) {
	return nil, errors.New("Failure!")
}

func (s *BadSerializable) Unserialize(bytes []byte) error {
	return errors.New("Failure!")
}

func validateBytes(t *testing.T, s Serializable) {
	gs, ok := s.(*GoodSerializable)
	if !ok {
		t.Fatalf("Cannot cast Serializable (%v) to GoodSerializable", s)
	}
	if !bytes.Equal(gs.ExpectedBytes, gs.actualBytes) {
		t.Errorf("Expected bytes (%s) does not match actual bytes (%s)", gs.ExpectedBytes, gs.actualBytes)
	}
}

func TestPersist(t *testing.T) {
	if _, err := os.Stat(tmpRoot); !os.IsNotExist(err) {
		t.Fatalf("Error: cannot run test when tmpRoot (%s) exists", tmpRoot)
	}
	defer os.RemoveAll(tmpRoot)

	data := []struct {
		name         string
		root         string
		context      string
		serializable Serializable
		success      bool
	}{
		{"BadRoot", "/this/file/doesnt/exist", tmpEmptyContext, &GoodSerializable{}, false},
		{"BadSerializable", tmpRoot, tmpEmptyContext, &BadSerializable{}, false},
		{"EmptyContext", tmpRoot, tmpEmptyContext, &GoodSerializable{}, true},
		{"SingletonContext",
			tmpRoot,
			tmpSingletonContext,
			&GoodSerializable{ExpectedBytes: []byte{'a', 'b', 'c'}},
			true},
		{"SingletonContextExists",
			tmpRoot,
			tmpSingletonContext,
			&GoodSerializable{ExpectedBytes: []byte{'c', 'b', 'a'}},
			true},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			t.Logf("Root: %s, Context: %s", datum.root, datum.context)

			p := Persister{datum.root}
			err := p.Persist(datum.context, datum.serializable)
			if !datum.success { // expecting failure...
				if err == nil {
					t.Errorf("Expected error, but got success")
				}
			} else { // expecting success...
				if err != nil {
					t.Errorf("Expected success, but got error: %s", err)
				} else {
					err := p.Unpersist(datum.context, datum.serializable)
					if err != nil {
						t.Errorf("Failed to unpersist persisted bytes: %s", err)
					}
					validateBytes(t, datum.serializable)
				}
			}
		})
	}
}

func TestUnpersist(t *testing.T) {
	data := []struct {
		name    string
		root    string
		context string
		s       Serializable
		success bool
	}{
		{"RootDoesntExist", badRoot, emptyContext, &GoodSerializable{}, false},
		{"ContextDoesntExist", goodRoot, badContext, &GoodSerializable{}, false},
		{"BadSerializable", goodRoot, emptyContext, &BadSerializable{}, false},
		{"EmptyContext", goodRoot, emptyContext, &GoodSerializable{ExpectedBytes: []byte{}}, true},
		{"SingletonContext",
			goodRoot,
			singletonContext,
			&GoodSerializable{ExpectedBytes: []byte{'a', 'b', 'c'}},
			true},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			t.Logf("Root: %s, Context: %s", datum.root, datum.context)

			p := Persister{datum.root}
			err := p.Unpersist(datum.context, datum.s)
			if !datum.success { // expecting failure...
				if err == nil {
					t.Errorf("Expected error, but got success")
				}
			} else { // expecting success...
				if err != nil {
					t.Errorf("Expected success, but got error: %s", err)
				} else {
					validateBytes(t, datum.s)
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
				err := p.Persist(datum.context, &GoodSerializable{})
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
