package storage

// This type represent an object that can marshal/unmarshal itself into bytes.
type Serializable interface {
	Serialize() ([]byte, error)
	Unserialize(bytes []byte) error
}
