package repo

type Object interface {
	Serialize() (string, error)
	Deserialize(data []byte) error
	GetFormat() string
}
