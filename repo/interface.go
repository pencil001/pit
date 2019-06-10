package repo

import (
	"compress/zlib"
	"fmt"
	"log"
	"path"

	"github.com/pencil001/pit/util"
)

type Object interface {
	Read(objSHA string) error
	Save() (string, error)
	Encode() ([]byte, error)
	Serialize() (string, error)
	Deserialize(data []byte) error
	GetFormat() string
}

type BaseObject struct {
	Object
	repo   *Repository
	format string
}

func (bo *BaseObject) Read(objSHA string) error {
	if bo.repo == nil {
		log.Panic("Repo mustn't be null")
	}
	format, content, err := bo.repo.parseObject(objSHA)
	if err != nil {
		return err
	}
	if format != bo.GetFormat() {
		return fmt.Errorf("Type is not correct: %v", format)
	}
	if err := bo.Deserialize([]byte(content)); err != nil {
		return err
	}
	return nil
}

func (bo *BaseObject) Save() (string, error) {
	if bo.repo == nil {
		log.Panic("Repo mustn't be null")
	}
	content, err := bo.Encode()
	if err != nil {
		return "", err
	}
	sha := util.CalcSHA(content)
	objDir := path.Join(bo.repo.gitDir, "objects", sha[:2])
	if err := util.CreateDir(objDir); err != nil {
		return "", err
	}
	fObj, err := util.CreateFile(path.Join(objDir, sha[2:]))
	if err != nil {
		return "", err
	}
	defer fObj.Close()

	w := zlib.NewWriter(fObj)
	defer w.Close()
	_, err = w.Write(content)
	if err != nil {
		return "", err
	}
	return sha, nil
}

func (bo *BaseObject) Encode() ([]byte, error) {
	str, err := bo.Serialize()
	if err != nil {
		return nil, err
	}
	content := fmt.Sprintf("%v %v\x00%v", bo.GetFormat(), len(str), str)
	return []byte(content), nil
}

func (bo *BaseObject) GetFormat() string {
	return bo.format
}
