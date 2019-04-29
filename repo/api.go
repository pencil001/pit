package repo

import (
	"log"
	"os"

	"github.com/pencil001/pit/util"
)

func Init(repoPath string) Repository {
	repo := createRepository(repoPath, true)

	isExist, err := util.IsExist(repo.workTree)
	if err != nil {
		log.Panicf("%v is not a directory!", repoPath)
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
