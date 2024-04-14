package git

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

type Git struct {
	worktree string
	gitdir   string
	repo     string
}

func (g *Git) CreateRepo() error {
	// .
	wf, err := os.Stat(g.worktree)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(g.worktree, 0755); err != nil {
				return fmt.Errorf("failed to create dir %v,  %v", g.worktree, err)
			}
		} else {
			return fmt.Errorf("failed to open worktree, %v", err)
		}
	} else {
		if !wf.IsDir() {
			return errors.New("worktree not a dir")
		}
	}

	// .git
	gf, err := os.Stat(g.gitdir)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(g.gitdir, 0755); err != nil {
				return fmt.Errorf("failed to create dir %v,  %v", g.gitdir, err)
			}
		} else {
			return fmt.Errorf("failed to open .git dir, %v", err)
		}
	} else {
		if !gf.IsDir() {
			return errors.New(".git is not a dir")
		}
	}

	// .git/branches
	// .git/objects
	// .git/refs/tags
	// .git/refs/heads
	if err := g.newGitDir("branches"); err != nil {
		return err
	}
	if err := g.newGitDir("objects"); err != nil {
		return err
	}
	if err := g.newGitDir("refs", "tags"); err != nil {
		return err
	}
	if err := g.newGitDir("refs", "heads"); err != nil {
		return err
	}

	// .git/description
	df, err := g.newGitFile("description")
	if err != nil {
		return fmt.Errorf("failed to create description file, %v", err)
	}
	df.WriteString("Unnamed repository; edit this file 'description' to name the repository.\n")
	df.Close()

	// .git/HEAD
	hf, err := g.newGitFile("HEAD")
	if err != nil {
		return fmt.Errorf("failed to create HEAD file, %v", err)
	}
	hf.WriteString("ref: refs/heads/master\n")
	hf.Close()

	// .git/config
	cf, err := g.newGitFile("config")
	if err != nil {
		return fmt.Errorf("failed to create config file, %v", err)
	}

	config := ini.Empty()
	coreSec, err := config.NewSection("core")
	if err != nil {
		return fmt.Errorf("failed to create config file, %v", err)
	}
	coreSec.NewKey("repositoryformatversion", "0")
	coreSec.NewKey("filemode", "false")
	coreSec.NewKey("bare", "false")
	config.WriteTo(cf)
	cf.Close()

	return nil
}

func (g *Git) init() {
	g.sanitizeWorktree()
	g.gitdir = g.worktree + string(os.PathSeparator) + ".git"
}

func (g *Git) newGitDir(p ...string) error {
	dpath := g.path(p...)
	dpath = g.gitdir + string(os.PathSeparator) + dpath
	return os.MkdirAll(dpath, 0755)
}

func (g *Git) newGitFile(p ...string) (*os.File, error) {
	dpath := g.path(p...)
	dpath = g.gitdir + string(os.PathSeparator) + dpath
	return os.Create(dpath)
}

func (g *Git) path(p ...string) string {
	sp := strings.Join(p, string(os.PathSeparator))
	return sp
}

func (g *Git) sanitizeWorktree() {
	rr := []rune(g.worktree)
	i := len(rr) - 1
	for ; i >= 0; i++ {
		if rr[i] != os.PathSeparator {
			break
		}
	}
	if i != len(rr)-1 {
		g.worktree = string(rr[:i+1])
	}
}

func NewGit(path string, repo string) *Git {
	git := Git{
		worktree: path,
		repo:     repo,
	}
	git.init()
	fmt.Printf("Repo: %v, Worktree: %s, Gitpath: %v\n", git.repo, git.worktree, git.gitdir)
	return &git
}
