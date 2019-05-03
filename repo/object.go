package repo

import (
	"fmt"
)

type Object interface {
	Serialize() []byte
	Deserialize(data []byte)
	GetFormat() string
	ToObjectBytes() []byte
}

type Blob struct {
	format string
	data   []byte
}

func createBlob(data []byte) *Blob {
	blob := Blob{
		format: TypeBlob,
	}
	if data != nil {
		blob.Deserialize(data)
	}
	return &blob
}

func (b *Blob) Serialize() []byte {
	return b.data
}

func (b *Blob) Deserialize(data []byte) {
	b.data = data
}

func (b *Blob) GetFormat() string {
	return b.format
}

func (b *Blob) ToObjectBytes() []byte {
	content := fmt.Sprintf("%v %v\x00%v", b.format, len(b.data), string(b.data))
	return []byte(content)
}
