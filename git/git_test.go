package git

import (
	"testing"
)

var (
	basePath = "../.."
)

func TestCreateRepo(t *testing.T) {
	git := NewGit(basePath+"/testgit", "myrepo")
	err := git.CreateRepo()
	if err != nil {
		t.Fatal(err)
	}
}
