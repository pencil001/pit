package repo

import (
	"fmt"
	"log"
	"path"

	"github.com/go-ini/ini"
	"github.com/pencil001/pit/util"
)

type Repository struct {
	workTree string
	gitDir   string
}

func createRepository(repoPath string, force bool) Repository {
	repo := Repository{
		workTree: repoPath,
		gitDir:   path.Join(repoPath, ".git"),
	}

	if !force {
		isDir, err := util.IsDir(repo.gitDir)
		if err != nil || !isDir {
			log.Panic(fmt.Sprintf("Not a Git repository %v", repoPath))
		}

		cfgFile := path.Join(repo.gitDir, "config")
		isExist, err := util.IsExist(cfgFile)
		if err != nil || !isExist {
			log.Panic("Configuration file missing")
		}

		f, err := ini.Load(cfgFile)
		if err != nil {
			log.Panic(err)
		}
		ver, err := f.Section("core").Key("repositoryformatversion").Int()
		if err != nil {
			log.Panic(fmt.Sprintf("Unanalyzable repositoryformatversion: %v", err))
		}
		if ver != 0 {
			log.Panic(fmt.Sprintf("Unsupported repositoryformatversion %v", ver))
		}
	}

	return repo
}

func (r Repository) initGitDir() error {
	if err := util.CreateDir(path.Join(r.gitDir, "branches")); err != nil {
		return err
	}
	if err := util.CreateDir(path.Join(r.gitDir, "objects")); err != nil {
		return err
	}
	if err := util.CreateDir(path.Join(r.gitDir, "refs", "tags")); err != nil {
		return err
	}
	if err := util.CreateDir(path.Join(r.gitDir, "refs", "heads")); err != nil {
		return err
	}

	fDesc, err := util.CreateFile(path.Join(r.gitDir, "description"))
	if err != nil {
		return err
	}
	defer fDesc.Close()
	fDesc.WriteString("Unnamed repository; edit this file 'description' to name the repository.\n")

	fHead, err := util.CreateFile(path.Join(r.gitDir, "HEAD"))
	if err != nil {
		return err
	}
	defer fHead.Close()
	fHead.WriteString("ref: refs/heads/master\n")

	cfgPath := path.Join(r.gitDir, "config")
	fConfig, err := util.CreateFile(cfgPath)
	if err != nil {
		return err
	}
	defer fConfig.Close()

	iniFile, err := ini.Load(cfgPath)
	if err != nil {
		return err
	}
	coreSect := iniFile.Section("core")
	coreSect.NewKey("repositoryformatversion", "0")
	coreSect.NewKey("filemode", "false")
	coreSect.NewKey("bare", "false")
	iniFile.WriteTo(fConfig)

	return nil
}
