package empty_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ajm113/git-trash/internal/actions/empty"
	"github.com/ajm113/git-trash/internal/git/mock"
	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	ctx := context.Background()
	trashPrefix := "git-trash/"

	t.Run("return error when fetching branches fail", func(t *testing.T) {
		git := &mock.MockGit{}
		git.On("Branches").Return(([]string)(nil), errors.New("unexpected error"))

		err := empty.All(ctx, git, trashPrefix)
		require.Error(t, err)
	})

	t.Run("delete single branch thats in the bin", func(t *testing.T) {
		git := &mock.MockGit{}
		trashBranch := trashPrefix + "test"
		git.On("Branches").Return([]string{trashBranch}, nil)
		git.On("DeleteBranch", trashBranch).Return(nil)

		err := empty.All(ctx, git, trashPrefix)
		require.NoError(t, err)
	})
	t.Run("delete single branch thats in the bin that has multiple branches", func(t *testing.T) {
		git := &mock.MockGit{}
		trashBranch := trashPrefix + "test"
		git.On("Branches").Return([]string{trashBranch, "test", "main"}, nil)
		git.On("DeleteBranch", trashBranch).Return(nil)

		err := empty.All(ctx, git, trashPrefix)
		require.NoError(t, err)
	})
	t.Run("delete multiple branches in the bin", func(t *testing.T) {
		git := &mock.MockGit{}
		trashBranch := trashPrefix + "test"
		trashBranch2 := trashPrefix + "test2"
		git.On("Branches").Return([]string{trashBranch, trashBranch2, "test", "main"}, nil)
		git.On("DeleteBranch", trashBranch).Return(nil)
		git.On("DeleteBranch", trashBranch2).Return(nil)

		err := empty.All(ctx, git, trashPrefix)
		require.NoError(t, err)
	})
	t.Run("returns error when attempting to delete a branch", func(t *testing.T) {
		git := &mock.MockGit{}
		trashBranch := trashPrefix + "test"
		git.On("Branches").Return([]string{trashBranch, "test", "main"}, nil)
		git.On("DeleteBranch", trashBranch).Return(errors.New("unexpected error"))

		err := empty.All(ctx, git, trashPrefix)
		require.Error(t, err)
	})
}
