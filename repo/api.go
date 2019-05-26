package repo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pencil001/pit/util"
)

const (
	TypeBlob   = "blob"
	TypeCommit = "commit"
	TypeTree   = "tree"
	TypeTag    = "tag"
)

func Init(repoPath string) *Repository {
	repo := createRepository(repoPath, true)

	isExist, err := util.IsExist(repo.workTree)
	if err != nil {
		log.Panicf("%v is not exist!", repoPath)
	}
	if isExist {
		if isDir, err := util.IsDir(repo.workTree); err != nil || !isDir {
			log.Panicf("%v is not a directory!", repoPath)
		}
	} else {
		if err := os.MkdirAll(repoPath, 0777); err != nil {
			log.Panicf("Create directory %v failed!", repoPath)
		}
	}

	if err := repo.initGitDir(); err != nil {
		log.Panic(err)
	}
	return repo
}

func Hash(filePath string, objType string, isStore bool) string {
	var repo *Repository
	if isStore {
		repo = createRepository(".", false)
	}

	isExist, err := util.IsExist(filePath)
	if err != nil || !isExist {
		log.Panicf("%v is not exist!", filePath)
	}
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}

	var obj Object
	switch objType {
	case TypeBlob:
		obj = createBlob(repo, fileData)
	case TypeCommit:
		obj = createCommit(repo, fileData)
	case TypeTree:
		obj = createTree(repo, fileData)
	default:
		log.Panicf("Unknown type: %v", objType)
	}

	if isStore {
		if err := obj.Save(); err != nil {
			log.Panic(err)
		}
	}

	bs, err := obj.Encode()
	if err != nil {
		log.Panic(err)
	}
	return util.CalcSHA(bs)
}

func Cat(objType string, objSHA string) string {
	repo := findRepo(".")

	var obj Object
	switch objType {
	case TypeBlob:
		obj = createBlob(repo, nil)
	case TypeCommit:
		obj = createCommit(repo, nil)
	case TypeTree:
		obj = createTree(repo, nil)
	default:
		log.Panicf("Unknown type: %v", objType)
	}

	err := obj.Read(objSHA)
	if err != nil {
		log.Panic(err)
	}

	bs, err := obj.Serialize()
	if err != nil {
		log.Panic(err)
	}
	return string(bs)
}

func Log(objSHA string) string {
	repo := findRepo(".")

	var sb strings.Builder
	sb.WriteString("digraph pit{\n")
	graphvizLog(&sb, repo, objSHA, map[string]bool{})
	sb.WriteString("}\n")
	return sb.String()
}

func ListTree(objSHA string) string {
	var err error

	repo := findRepo(".")
	obj, err := repo.readObject(objSHA)
	if err != nil {
		log.Panic(err)
	}

	format := obj.GetFormat()
	if format != TypeCommit && format != TypeTree {
		log.Panic("not a tree object")
	}

	var str string
	if format == TypeTree {
		tree := obj.(*Tree)
		str, err = tree.Display()
		if err != nil {
			log.Panic(err)
		}
	}
	if format == TypeCommit {
		commit := obj.(*Commit)
		for _, kv := range commit.kvlm {
			if kv.key == "tree" {
				treeSHA := kv.list[0]
				tree := createTree(repo, nil)
				if err := tree.Read(treeSHA); err != nil {
					log.Panic(err)
				}
				str, err = tree.Display()
				if err != nil {
					log.Panic(err)
				}
				break
			}
		}
	}
	return str
}

func findRepo(repoPath string) *Repository {
	gitPath := path.Join(repoPath, ".git")
	isDir, _ := util.IsDir(gitPath)
	if isDir {
		return createRepository(repoPath, false)
	}
	parentPath := path.Join(repoPath, "..")

	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		log.Panic(err)
	}
	absParentPath, err := filepath.Abs(parentPath)
	if err != nil {
		log.Panic(err)
	}

	if absParentPath == absRepoPath {
		log.Panic("No git directory.")
	}
	return findRepo(parentPath)
}

func graphvizLog(sb *strings.Builder, repo *Repository, sha string, seen map[string]bool) {
	if _, ok := seen[sha]; ok {
		return
	}
	seen[sha] = true

	commit := createCommit(repo, nil)
	err := commit.Read(sha)
	if err != nil {
		log.Panic(err)
	}

	isInit := true
	var parentValue []string
	for _, kl := range commit.kvlm {
		if kl.key == "parent" {
			isInit = false
			parentValue = kl.list
			break
		}
	}
	// Base case: the initial commit.
	if isInit {
		return
	}

	for _, v := range parentValue {
		sb.WriteString(fmt.Sprintf("c_%v -> c_%v;\n", sha, v))
		graphvizLog(sb, repo, v, seen)
	}
}
