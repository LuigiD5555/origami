package environment

import (
	"os/exec"
	"runtime"
	"strings"
)

type Snapshot struct {
	Schema          string `json:"schema"`
	GoVersion       string `json:"go_version"`
	OS              string `json:"os"`
	Arch            string `json:"arch"`
	GitCommit       string `json:"git_commit"`
	GitState        string `json:"git_state"`
	DirtyWorktree   bool   `json:"dirty_worktree"`
	ProtocolVersion string `json:"protocol_version"`
	BaselineVersion string `json:"baseline_version"`
}

func Capture(repo string) Snapshot {
	commit, commitOK := command(repo, "git", "rev-parse", "HEAD")
	status, statusOK := command(repo, "git", "status", "--porcelain")
	gitState := "UNAVAILABLE"
	dirty := false
	if commitOK && statusOK {
		if status == "" {
			gitState = "CLEAN"
		} else {
			gitState = "DIRTY"
			dirty = true
		}
	} else {
		commit = "UNKNOWN"
	}
	return Snapshot{Schema: "ohf.environment.v2", GoVersion: runtime.Version(), OS: runtime.GOOS, Arch: runtime.GOARCH, GitCommit: commit, GitState: gitState, DirtyWorktree: dirty, ProtocolVersion: "R3.10-LAB", BaselineVersion: "R3.10"}
}

func command(dir, name string, args ...string) (string, bool) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	b, err := cmd.Output()
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(b)), true
}
