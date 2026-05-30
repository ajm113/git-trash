package empty

import (
	"context"
	"fmt"
	"strings"

	"github.com/ajm113/git-trash/internal/git"
)

func All(ctx context.Context, git git.Git, trashPrefix string) error {
	branches, err := git.Branches()
	if err != nil {
		fmt.Printf("error: unable to get branches: %s\n", err)
		return err
	}

	trashBranches := []string{}
	for _, branch := range branches {
		if strings.HasPrefix(branch, trashPrefix) {
			trashBranches = append(trashBranches, branch)
		}
	}

	for _, branch := range trashBranches {
		err := git.DeleteBranch(branch)
		if err != nil {
			fmt.Printf("error: [%s] unable to delete branch: %s\n", branch, err)
			return err
		}

		fmt.Printf("[%s] deleted from local repo\n", branch)
	}

	return nil
}
