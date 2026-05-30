package mock

import (
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

func (m *MockGit) MoveBranch(oldName, newName string) error {
	args := m.Called(oldName, newName)
	return args.Error(0)
}

func (m *MockGit) DeleteBranch(branch string) error {
	args := m.Called(branch)
	return args.Error(0)
}

func (m *MockGit) LastCommitTimestamp(branch string) (time.Time, error) {
	args := m.Called(branch)
	return args.Get(0).(time.Time), args.Error(1)
}
