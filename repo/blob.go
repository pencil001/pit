package repo

import (
	"log"
)

type Blob struct {
	BaseObject
	data []byte
}

func createBlob(repo *Repository, data []byte) *Blob {
	blob := &Blob{
		BaseObject: BaseObject{
			repo:   repo,
			format: TypeBlob,
		},
		data: []byte{},
	}
	blob.BaseObject.Object = blob
	if data != nil {
		err := blob.Deserialize(data)
		if err != nil {
			log.Panic(err)
		}
	}
	return blob
}

func (b *Blob) Serialize() (string, error) {
	return string(b.data), nil
}

func (b *Blob) Deserialize(data []byte) error {
	b.data = data
	return nil
}
