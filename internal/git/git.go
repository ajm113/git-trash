package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Git interface {
	Branches() ([]string, error)
	MoveBranch(string, string) error
	DeleteBranch(string) error
	LastCommitTimestamp(string) (time.Time, error)
}

type client struct{}

func New() Git {
	return &client{}
}

func (c *client) Branches() ([]string, error) {
	out, err := exec.Command("git", "branch").Output()
	if err != nil {
		return nil, fmt.Errorf("git branch failed: %w", err)
	}

	branches := []string{}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "* ") {
			current := strings.TrimPrefix(line, "* ")
			branches = append(branches, current)
		} else {
			branches = append(branches, line)
		}
	}

	return branches, nil
}

func (c *client) MoveBranch(oldName, newName string) error {
	_, err := exec.Command("git", "branch", "-m", oldName, newName).Output()
	if err != nil {
		return fmt.Errorf("git branch -m failed: %w", err)
	}

	return nil
}

func (c *client) DeleteBranch(branch string) error {
	_, err := exec.Command("git", "branch", "-D", branch).Output()
	if err != nil {
		return fmt.Errorf("git branch -D failed: %w", err)
	}

	return nil
}

func (c *client) LastCommitTimestamp(branch string) (time.Time, error) {
	out, err := exec.Command("git", "log", fmt.Sprintf("--max-count=%d", 1), "--format=%ct", branch).Output()
	if err != nil {
		return time.Time{}, fmt.Errorf("git log failed: %w", err)
	}

	unix, err := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing commit timestamp %q: %w", out, err)
	}

	return time.Unix(unix, 0), nil
}
