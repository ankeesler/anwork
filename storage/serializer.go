package storage

type Serializer interface {
	Unserialize(bytes []byte) (interface{}, error)
	Serialize(interface{}) ([]byte, error)
}
