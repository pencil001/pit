package repo

import (
	"log"
)

type Blob struct {
	format string
	data   []byte
}

func createBlob(data []byte) *Blob {
	blob := Blob{
		format: TypeBlob,
	}
	if data != nil {
		err := blob.Deserialize(data)
		if err != nil {
			log.Panic(err)
		}
	}
	return &blob
}

func (b *Blob) Serialize() (string, error) {
	return string(b.data), nil
}

func (b *Blob) Deserialize(data []byte) error {
	b.data = data
	return nil
}

func (b *Blob) GetFormat() string {
	return b.format
}
