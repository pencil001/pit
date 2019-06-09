package repo

type Tag struct {
	*Commit
}

func createTag(repo *Repository, data []byte) *Tag {
	tag := &Tag{
		Commit: createCommit(repo, data),
	}
	tag.format = TypeTag
	return tag
}
