package mock

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockGit struct {
	mock.Mock
}

func (m *MockGit) Branches() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockGit) MoveBranches(oldName, newName string) error {
	args := m.Called(oldName, newName)
	return args.Error(0)
}

func (c *client) DeleteBranch(branch string) error {
	out, err := exec.Command("git", "branch", "-D", branch).Output()
	if err != nil {
		return fmt.Errorf("git branch -D failed: %w", err)
	}

	return nil
}

func (c *client) LastCommitTimestamp(branch string) (time.Time, error) {
	out, err := exec.Command("git", "log", fmt.Sprintf("--max-count=%d", limit), "--format=\"%ct\"").Output()
	if err != nil {
		return "", fmt.Errorf("git branch -m failed: %w", err)
	}

	return string(out), nil
}
