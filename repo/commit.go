package repo

import (
	"fmt"
	"log"
	"strings"

	"github.com/pencil001/pit/util"
)

type KList struct {
	key  string
	list []string
}

type Commit struct {
	BaseObject
	kvlm []KList
}

func createCommit(repo *Repository, data []byte) *Commit {
	commit := &Commit{
		BaseObject: BaseObject{
			repo:   repo,
			format: TypeCommit,
		},
		kvlm: []KList{},
	}
	commit.BaseObject.Object = commit
	if data != nil {
		err := commit.Deserialize(data)
		if err != nil {
			log.Panic(err)
		}
	}
	return commit
}

func (c *Commit) Serialize() (string, error) {
	m := ""
	var sb strings.Builder
	for _, kl := range c.kvlm {
		if kl.key != "" {
			for _, s := range kl.list {
				fs := strings.ReplaceAll(s, "\n", "\n ")
				sb.WriteString(fmt.Sprintf("%v %v\n", kl.key, fs))
			}
		} else {
			if len(kl.list) == 1 {
				m = kl.list[0]
			}
		}
	}
	sb.WriteString(fmt.Sprintf("\n%v", m))
	return sb.String(), nil
}

func (c *Commit) Deserialize(data []byte) error {
	return c.parse([]rune(string(data)))
}

func (c *Commit) parse(rs []rune) error {
	idxSpace := util.FindInRunes(rs, ' ', 0)
	idxNewLine := util.FindInRunes(rs, '\n', 0)

	// Base case
	// =========
	// If newline appears first (or there's no space at all, in which
	// case find returns -1), we assume a blank line.  A blank line
	// means the remainder of the data is the message.
	if idxSpace < 0 || idxNewLine < idxSpace {
		if idxNewLine != 0 {
			return fmt.Errorf("No blank line")
		}
		c.kvlm = append(c.kvlm, KList{
			key:  "",
			list: []string{string(rs[idxNewLine+1:])},
		})
		return nil
	}

	// Recursive case
	// ==============
	// we read a key-value pair and recurse for the next.
	key := string(rs[:idxSpace])

	// Find the end of the value.  Continuation lines begin with a
	// space, so we loop until we find a "\n" not followed by a space.
	end := 0
	for {
		end = util.FindInRunes(rs, '\n', end+1)
		if rs[end+1] != ' ' {
			break
		}
	}

	// Grab the value
	// Also, drop the leading space on continuation lines
	value := strings.ReplaceAll(string(rs[idxSpace+1:end]), "\n ", "\n")
	exist := false
	for _, kl := range c.kvlm {
		if kl.key == key {
			kl.list = append(kl.list, value)
			exist = true
			break
		}
	}
	if !exist {
		c.kvlm = append(c.kvlm, KList{
			key:  key,
			list: []string{value},
		})
	}
	return c.parse(rs[end+1:])
}
