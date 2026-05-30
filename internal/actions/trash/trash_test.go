package trash_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ajm113/git-trash/internal/actions/trash"
	"github.com/ajm113/git-trash/internal/git/mock"
	"github.com/stretchr/testify/require"
)

func TestByMatch(t *testing.T) {
	ctx := context.Background()
	trashPrefix := "git-trash/"
	protectedBranches := []string{"main", "master"}

	t.Run("return error when no match string is provided", func(t *testing.T) {
		git := &mock.MockGit{}
		err := trash.ByMatch(ctx, git, "", protectedBranches, trashPrefix)
		require.Error(t, err)
	})
	t.Run("return error when fetching branches fail", func(t *testing.T) {
		git := &mock.MockGit{}
		git.On("Branches").Return(([]string)(nil), errors.New("unexpected error"))

		err := trash.ByMatch(ctx, git, "test", protectedBranches, trashPrefix)
		require.Error(t, err)
	})
	t.Run("move single branch to trash", func(t *testing.T) {
		git := &mock.MockGit{}
		testBranch := "test"
		testTrashBranch := trashPrefix + "test"

		git.On("Branches").Return([]string{testBranch}, nil)
		git.On("MoveBranch", testBranch, testTrashBranch).Return(nil)

		err := trash.ByMatch(ctx, git, "test", protectedBranches, trashPrefix)
		require.NoError(t, err)
	})
	t.Run("move single branch to trash returns error", func(t *testing.T) {
		git := &mock.MockGit{}
		testBranch := "test"
		testTrashBranch := trashPrefix + "test"

		git.On("Branches").Return([]string{testBranch}, nil)
		git.On("MoveBranch", testBranch, testTrashBranch).Return(errors.New("unexpected error"))

		err := trash.ByMatch(ctx, git, "test", protectedBranches, trashPrefix)
		require.Error(t, err)
	})
	t.Run("noop: user attempts to move main to trash", func(t *testing.T) {
		git := &mock.MockGit{}
		testBranch := "main"

		git.On("Branches").Return([]string{testBranch}, nil)

		err := trash.ByMatch(ctx, git, "main", protectedBranches, trashPrefix)
		require.NoError(t, err)
	})
	t.Run("noop: user attempts to move a protected branch to trash by using glob", func(t *testing.T) {
		git := &mock.MockGit{}
		git.On("Branches").Return([]string{"test", "main", "master"}, nil)

		err := trash.ByMatch(ctx, git, "m*", protectedBranches, trashPrefix)
		require.NoError(t, err)
	})
	t.Run("glob match test branch with protected", func(t *testing.T) {
		git := &mock.MockGit{}
		testBranch := "maw"
		testTrashBranch := trashPrefix + "maw"

		git.On("Branches").Return([]string{testBranch, "main", "master"}, nil)
		git.On("MoveBranch", testBranch, testTrashBranch).Return(nil)

		err := trash.ByMatch(ctx, git, "m*", protectedBranches, trashPrefix)
		require.NoError(t, err)
	})
	t.Run("noop: skip branches that already exist in trash bin", func(t *testing.T) {
		git := &mock.MockGit{}
		testBranch := "maw"
		testTrashBranch := trashPrefix + "maw"

		git.On("Branches").Return([]string{testBranch, testTrashBranch, "main", "master"}, nil)

		err := trash.ByMatch(ctx, git, "maw", protectedBranches, trashPrefix)
		require.NoError(t, err)
	})
}

func TestByDays(t *testing.T) {
	ctx := context.Background()
	trashPrefix := "git-trash/"
	protectedBranches := []string{"main", "master"}

	t.Run("return error when the input for days is 0", func(t *testing.T) {
		git := &mock.MockGit{}

		err := trash.ByDays(ctx, git, 0, protectedBranches, trashPrefix)
		require.Error(t, err)
	})
	t.Run("return error when fetching branches fail", func(t *testing.T) {
		git := &mock.MockGit{}
		git.On("Branches").Return(([]string)(nil), errors.New("unexpected error"))

		err := trash.ByDays(ctx, git, 1, protectedBranches, trashPrefix)
		require.Error(t, err)
	})
	t.Run("move single branch to trash", func(t *testing.T) {
		git := &mock.MockGit{}
		testBranch := "test"
		testTrashBranch := trashPrefix + "test"

		git.On("Branches").Return([]string{testBranch, "main", "master"}, nil)
		git.On("LastCommitTimestamp", testBranch).Return(time.Now().AddDate(0, 0, 2), nil)
		git.On("MoveBranch", testBranch, testTrashBranch).Return(nil)

		err := trash.ByDays(ctx, git, 1, protectedBranches, trashPrefix)
		require.NoError(t, err)
	})
	t.Run("noop: dont move young branches to trash", func(t *testing.T) {
		git := &mock.MockGit{}
		testBranch := "test"

		git.On("Branches").Return([]string{testBranch, "main", "master"}, nil)
		git.On("LastCommitTimestamp", testBranch).Return(time.Now().AddDate(0, 0, 2), nil)

		err := trash.ByDays(ctx, git, 360, protectedBranches, trashPrefix)
		require.NoError(t, err)
	})
	t.Run("noop: dont move young branches to trash and ignore trash bin", func(t *testing.T) {
		git := &mock.MockGit{}
		testBranch := "test"

		git.On("Branches").Return([]string{testBranch, trashPrefix + "abc", "main", "master"}, nil)
		git.On("LastCommitTimestamp", testBranch).Return(time.Now().AddDate(0, 0, 2), nil)

		err := trash.ByDays(ctx, git, 360, protectedBranches, trashPrefix)
		require.NoError(t, err)
	})
	t.Run("noop: skip branches that already exist in trash bin", func(t *testing.T) {
		git := &mock.MockGit{}
		testBranch := "maw"
		testTrashBranch := trashPrefix + "maw"

		git.On("Branches").Return([]string{testBranch, testTrashBranch, "main", "master"}, nil)
		git.On("LastCommitTimestamp", testBranch).Return(time.Now().AddDate(0, 0, 2), nil)

		err := trash.ByDays(ctx, git, 1, protectedBranches, trashPrefix)
		require.NoError(t, err)
	})
}
