package trash

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/ajm113/git-trash/internal/git"
	"github.com/ajm113/git-trash/internal/globparser"
)

func ByMatch(ctx context.Context, git git.Git, match string, protectedBranches []string, trashPrefix string) error {
	if match == "" {
		fmt.Println("error:" + ErrNoBranchOrMatchProvided.Error())
		return ErrNoBranchOrMatchProvided
	}

	branches, err := git.Branches()
	if err != nil {
		fmt.Println("error:" + err.Error())
		return err
	}

	regex, err := globparser.Regex(match)
	if err != nil {
		fmt.Println("error:" + err.Error())
		return err
	}

	branchesToTrash := []string{}
	for b := range slices.Values(branches) {
		if trashPrefix != "" && strings.HasPrefix(b, trashPrefix) {
			continue
		}

		if slices.Contains(protectedBranches, b) {
			continue
		}

		if regex.Match([]byte(b)) {
			branchesToTrash = append(branchesToTrash, b)
		}
	}

	return moveBranchesToTrash(git, branches, branchesToTrash, trashPrefix)
}

func ByDays(ctx context.Context, git git.Git, days int, protectedBranches []string, trashPrefix string) error {
	if days <= 0 {
		err := ErrInvalidDayCount
		fmt.Println("error: " + err.Error())
		return err
	}

	branches, err := git.Branches()
	if err != nil {
		fmt.Println("error:" + err.Error())
		return err
	}

	cutoffTimestamp := time.Now().AddDate(0, 0, days*-1).UTC()

	branchesToTrash := []string{}
	for b := range slices.Values(branches) {
		if trashPrefix != "" && strings.HasPrefix(b, trashPrefix) {
			continue
		}

		if slices.Contains(protectedBranches, b) {
			continue
		}

		lastCommitTimestamp, err := git.LastCommitTimestamp(b)
		if err != nil {
			fmt.Println("error: failed getting last commit: " + err.Error())
			return err
		}

		if lastCommitTimestamp.Before(cutoffTimestamp) {
			branchesToTrash = append(branchesToTrash, b)
		}
	}

	return moveBranchesToTrash(git, branches, branchesToTrash, trashPrefix)
}

func moveBranchesToTrash(git git.Git, currentBranches, branchesToTrash []string, trashPrefix string) error {
	if len(branchesToTrash) == 0 {
		fmt.Println("warn: no branches to move to trash")
		return nil
	}

	for b := range slices.Values(branchesToTrash) {

		// make sure we won't have a collision first.
		if slices.Contains(currentBranches, trashPrefix+b) {
			fmt.Printf("warn: [%s] would result in a conflict if moved to trash bin. Please rename branch with trash prefix or change the current branch name!\n", b)
			continue
		}

		err := git.MoveBranch(b, trashPrefix+b)
		if err != nil {
			fmt.Printf("error: [%s] failed moving branch: %s\n", b, err.Error())
			return err
		}
		fmt.Printf("[%s] moved to trash bin -> %s\n", b, trashPrefix+b)
	}

	return nil
}
