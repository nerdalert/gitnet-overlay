package control

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

// Git commands
type GitCmd struct {
	Dir string
}

// NewGit create
func newGit(dir string) *GitCmd {
	return &GitCmd{Dir: dir}
}

// Update method executes 'git pull'
func (g *GitCmd) update() (cmd *exec.Cmd) {
	args := []string{"pull"}
	cmd = gitCmd(args)
	cmd.Dir = g.Dir
	return
}

// Fetch unused but added in case 'git fetch origin'
func (g *GitCmd) fetch() (cmd *exec.Cmd) {
	args := []string{"fetch", "origin"}
	cmd = gitCmd(args)
	cmd.Dir = g.Dir
	return
}

// updateCurrent update the current branch
func (g *GitCmd) updateCurrent() (cmd *exec.Cmd) {
	args := []string{"pull", "origin", currentBranch(g.Dir)}
	cmd = gitCmd(args)
	cmd.Dir = g.Dir
	return
}

func (g *GitCmd) clone(repo string) (cmd *exec.Cmd) {
	args := []string{"clone", repo, g.Dir}
	cmd = gitCmd(args)
	return
}

func currentBranch(path string) string {
	args := []string{"rev-parse", "--abbrev-ref", "HEAD"}
	cmd := exec.Command("git", args...)
	output := new(bytes.Buffer)
	cmd.Stdout = output
	cmd.Stderr = output
	cmd.Dir = path
	err := cmd.Run()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(output.String())
}

func gitCmd(args []string) (cmd *exec.Cmd) {
	cmd = exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return
}

func debugGitCmd(outs []byte) {
	if len(outs) > 0 {
		log.Debugf("Git output: %s\n", string(outs))
	}
}
