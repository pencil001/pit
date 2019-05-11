package repo

import (
	"fmt"
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

func (b *Blob) Serialize() ([]byte, error) {
	return b.data, nil
}

func (b *Blob) Deserialize(data []byte) error {
	b.data = data
	return nil
}

func (b *Blob) GetFormat() string {
	return b.format
}

func (b *Blob) ToObjectBytes() ([]byte, error) {
	content := fmt.Sprintf("%v %v\x00%v", b.format, len(b.data), string(b.data))
	return []byte(content), nil
}
