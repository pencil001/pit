package repo

import (
	"io/ioutil"
	"log"
	"os"
	"path"

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
		obj = createBlob(fileData)
	default:
		log.Panicf("Unknown type: %v", objType)
	}

	if isStore {
		err := repo.saveObject(obj)
		if err != nil {
			log.Panic(err)
		}
	}

	return util.CalcSHA(obj.ToObjectBytes())
}

func Cat(objType string, objSHA string) string {
	repo := findRepo(".")

	var obj Object
	switch objType {
	case TypeBlob:
		obj = createBlob(nil)
	default:
		log.Panicf("Unknown type: %v", objType)
	}

	err := repo.readObject(objSHA, obj)
	if err != nil {
		log.Panic(err)
	}
	return string(obj.Serialize())
}

func findRepo(repoPath string) *Repository {
	gitPath := path.Join(repoPath, ".git")
	isDir, _ := util.IsDir(gitPath)
	if isDir {
		return createRepository(repoPath, false)
	}

	parentPath := path.Base(repoPath)
	if parentPath == repoPath {
		log.Panic("No git directory.")
	}
	return findRepo(parentPath)
}
