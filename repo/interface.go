package repo

type Object interface {
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
	GetFormat() string
	ToObjectBytes() ([]byte, error)
}
