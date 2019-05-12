package repo

import (
	"fmt"
	"path"
	"path/filepath"
	"testing"
)

func TestPath(t *testing.T) {
	pth := "."
	parent, err := filepath.Abs(path.Join(pth, "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(parent)
}
