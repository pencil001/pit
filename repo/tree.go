package repo

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pencil001/pit/util"
)

type Leaf struct {
	mode string
	path string
	sha  string
}

type Tree struct {
	BaseObject
	leaves []Leaf
}

func createTree(repo *Repository, data []byte) *Tree {
	tree := &Tree{
		BaseObject: BaseObject{
			repo:   repo,
			format: TypeTree,
		},
		leaves: []Leaf{},
	}
	tree.BaseObject.Object = tree
	if data != nil {
		err := tree.Deserialize(data)
		if err != nil {
			log.Panic(err)
		}
	}
	return tree
}

func (t *Tree) Display() (string, error) {
	var sb strings.Builder
	for _, leaf := range t.leaves {
		obj, err := t.BaseObject.repo.readObject(leaf.sha)
		if err != nil {
			log.Panic(err)
		}
		sb.WriteString(fmt.Sprintf("%v %v %v\t%v\n", leaf.mode, obj.GetFormat(), leaf.sha, leaf.path))
	}
	return sb.String(), nil
}

func (t *Tree) Serialize() (string, error) {
	var sb strings.Builder
	for _, leaf := range t.leaves {
		sb.WriteString(fmt.Sprintf("%v %v\x00%v", leaf.mode, leaf.path, string(util.HexStrToBytes(leaf.sha))))
	}
	return sb.String(), nil
}

func (t *Tree) Deserialize(data []byte) error {
	pos := 0
	max := len(data)
	for pos < max {
		var leaf Leaf
		pos, leaf = t.parseOneLeaf(data, pos)
		t.leaves = append(t.leaves, leaf)
	}
	return nil
}

func (t *Tree) parseOneLeaf(rs []byte, start int) (int, Leaf) {
	idxSpace := util.FindInBytes(rs, ' ', start)
	im, err := strconv.ParseUint(string(rs[start:idxSpace]), 10, 32)
	if err != nil {
		log.Panic(err)
	}
	mode := fmt.Sprintf("%06d", im)

	idxNULL := util.FindInBytes(rs, '\x00', idxSpace)
	path := string(rs[idxSpace+1 : idxNULL])

	sha := util.BytesToHexStr(rs[idxNULL+1 : idxNULL+21])
	return idxNULL + 21, Leaf{mode, path, sha}
}
