package repo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
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
	case TypeTag:
		obj = createTag(repo, fileData)
	default:
		log.Panicf("Unknown type: %v", objType)
	}

	sha := ""
	if isStore {
		sha, err = obj.Save()
		if err != nil {
			log.Panic(err)
		}
	} else {
		bs, err := obj.Encode()
		if err != nil {
			log.Panic(err)
		}
		sha = util.CalcSHA(bs)
	}
	return sha
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
	case TypeTag:
		obj = createTag(repo, nil)
	default:
		log.Panicf("Unknown type: %v", objType)
	}

	objSHA = RevParse(objSHA, objType)
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
		tree := getTreeByCommit(commit)
		if tree == nil {
			log.Panic("No tree in commit")
		}
		str, err = tree.Display()
		if err != nil {
			log.Panic(err)
		}
	}
	return str
}

func Checkout(objSHA string, dir string) {
	repo := findRepo(".")
	obj, err := repo.readObject(objSHA)
	if err != nil {
		log.Panic(err)
	}

	format := obj.GetFormat()
	if format == TypeCommit {
		commit := obj.(*Commit)
		obj = getTreeByCommit(commit)
		format = obj.GetFormat()
		if obj == nil {
			log.Panic("No tree in commit")
		}
	}

	if err := ensureEmptyDir(dir); err != nil {
		log.Panic(err)
	}

	switch format {
	case TypeTree:
		if err := checkoutTree(obj, dir); err != nil {
			log.Panic(err)
		}
	}
}

func checkoutTree(treeObj Object, dir string) error {
	tree := treeObj.(*Tree)
	for _, leaf := range tree.leaves {
		subObj, err := tree.repo.readObject(leaf.sha)
		if err != nil {
			return err
		}

		destPath := path.Join(dir, leaf.path)
		switch subObj.GetFormat() {
		case TypeBlob:
			if err := checkoutBlob(subObj, destPath); err != nil {
				return err
			}
		case TypeTree:
			if err := ensureEmptyDir(destPath); err != nil {
				return err
			}
			if err := checkoutTree(subObj, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkoutBlob(blobObj Object, filePath string) error {
	fBlob, err := util.CreateFile(filePath)
	if err != nil {
		return err
	}
	defer fBlob.Close()

	content, err := blobObj.Serialize()
	if err != nil {
		return err
	}
	fBlob.WriteString(content)
	return nil
}

func ensureEmptyDir(dir string) error {
	isExist, err := util.IsExist(dir)
	if err != nil {
		return err
	}
	if !isExist {
		if err := util.CreateDir(dir); err != nil {
			return err
		}
	}

	isEmpty, err := util.IsEmptyDir(dir)
	if err != nil {
		return err
	}
	if !isEmpty {
		return errors.New("This dir is not empty")
	}
	return nil
}

func getTreeByCommit(commit *Commit) *Tree {
	for _, kv := range commit.kvlm {
		if kv.key == "tree" {
			treeSHA := kv.list[0]
			tree := createTree(commit.repo, nil)
			if err := tree.Read(treeSHA); err != nil {
				log.Panic(err)
			}
			return tree
		}
	}
	return nil
}

func ShowRefs(prefix string, withHash bool) string {
	repo := findRepo(".")

	refs, err := repo.getRefs()
	if err != nil {
		log.Panic(err)
	}

	// TODO: sort key
	var sb strings.Builder
	for k, v := range refs {
		if strings.HasPrefix(k, prefix) {
			k = strings.TrimPrefix(k, prefix)
			if withHash {
				sb.WriteString(fmt.Sprintf("%v %v\n", v, k))
			} else {
				sb.WriteString(fmt.Sprintf("%v\n", k))
			}
		}
	}
	return sb.String()
}

func ShowOrNewTag(tagName string, objSHA string, isCreateObject bool) string {
	if tagName == "" {
		return ShowRefs("refs/tags/", false)
	}

	// TODO: add tag
	repo := findRepo(".")
	if !isCreateObject {
		if err := repo.writeRef(path.Join("refs", "tags", tagName), objSHA); err != nil {
			log.Panic(err)
		}
	} else {
		tag := createTag(repo, nil)
		tag.kvlm = []KList{
			KList{key: "object", list: []string{objSHA}},
			KList{key: "type", list: []string{TypeCommit}},
			KList{key: "tag", list: []string{"vincent"}},
			KList{key: "tagger", list: []string{"vincent <pencil.xu@gmail.com>"}},
			KList{key: "", list: []string{"This is the commit message that should have come from the user\n"}},
		}

		sha, err := tag.Save()
		if err != nil {
			log.Panic(err)
		}

		if err := repo.writeRef(path.Join("refs", "tags", tagName), sha); err != nil {
			log.Panic(err)
		}
	}
	return ""
}

func RevParse(objRev, revType string) string {
	repo := findRepo(".")
	candidates, err := resolveObjectRev(repo, objRev)
	if err != nil {
		log.Panic(err)
	}

	if len(candidates) == 0 {
		log.Panicf("No such reference %v.", objRev)
	}
	if len(candidates) > 1 {
		log.Panicf("Ambiguous reference %v: Candidates are:\n%v.", objRev, strings.Join(candidates, "\n"))
	}

	hash := candidates[0]
	for {
		obj, err := repo.readObject(hash)
		if err != nil {
			log.Panic(err)
		}

		if revType != "" && obj.GetFormat() == revType {
			break
		}

		if obj.GetFormat() == TypeTag {
			isFollowed := false
			tag := obj.(*Tag)
			for _, kv := range tag.kvlm {
				if kv.key == "object" {
					hash = kv.list[0]
					isFollowed = true
					break
				}
			}
			if isFollowed {
				continue
			}
		}
		if obj.GetFormat() == TypeCommit && revType == TypeTree {
			isFollowed := false
			commit := obj.(*Commit)
			for _, kv := range commit.kvlm {
				if kv.key == "tree" {
					hash = kv.list[0]
					isFollowed = true
					break
				}
			}
			if isFollowed {
				continue
			}
		}

		if revType != "" && obj.GetFormat() != revType {
			log.Panic("wrong rev type")
		}
		break
	}
	return hash
}

func resolveObjectRev(repo *Repository, objRev string) ([]string, error) {
	if objRev == "" {
		return []string{}, nil
	}

	if objRev == "HEAD" {
		hash, err := repo.readRef(objRev, make(map[string]string))
		if err != nil {
			return nil, err
		}
		return []string{hash}, nil
	}

	if objRev == "master" {
		objRev = path.Join("refs", "heads", "master")
	}

	split := func(c rune) bool {
		return c == '/' || c == '\\'
	}
	inputSegments := strings.FieldsFunc(objRev, split)

	refs, err := repo.getRefs()
	if err != nil {
		return nil, err
	}
	candidates := []string{}
	for key, hash := range refs {
		keySegments := strings.FieldsFunc(key, split)
		if matchSegments(inputSegments, keySegments) {
			candidates = append(candidates, hash)
		}
	}

	hashRegex := regexp.MustCompile(`^[0-9A-Fa-f]{1,16}$`)
	if hashRegex.MatchString(objRev) {
		objRev = strings.ToLower(objRev)
		if len(objRev) == 40 {
			candidates = append(candidates, objRev)
		} else if len(objRev) >= 4 {
			prefix := objRev[:2]
			dir := path.Join(repo.gitDir, "objects", prefix)
			entries, err := ioutil.ReadDir(dir)
			if err != nil {
				return nil, err
			}
			for _, e := range entries {
				if strings.HasPrefix(e.Name(), objRev[2:]) {
					candidates = append(candidates, prefix+e.Name())
				}
			}
		}
	}
	return candidates, nil
}

func matchSegments(srcSegments, dstSegments []string) bool {
	i, j := len(srcSegments)-1, len(dstSegments)-1
	for i >= 0 && j >= 0 {
		if srcSegments[i] != dstSegments[j] {
			return false
		}
		i, j = i-1, j-1
	}
	if i >= 0 {
		return false
	}
	return true
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
