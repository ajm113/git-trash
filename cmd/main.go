package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ajm113/git-trash/internal/actions/empty"
	"github.com/ajm113/git-trash/internal/actions/trash"
	"github.com/ajm113/git-trash/internal/git"
	"github.com/urfave/cli/v3"
)

var (
	Version           = "dev"
	Commit            = "unknown"
	Date              = "unknown"
	TrashPrefix       = "git-trash/"
	ProtectedBranches = []string{"main", "master"}
)

func main() {

	git := git.New()

	cmd := &cli.Command{
		Name:  "git-trash",
		Usage: "Moves branches to a trash bin by simply append a prefix that marks them for deletion",
		Commands: []*cli.Command{
			{
				Name:    "trash",
				Aliases: []string{"t"},
				Usage:   "renames branch to trash prefix by matching the branch name",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return trash.ByMatch(ctx, git, getLastArg(cmd.Args()), ProtectedBranches, TrashPrefix)
				},
			},
			{
				Name:    "days",
				Aliases: []string{"d"},
				Usage:   "moves branches to trash bin by number of days",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					days, err := strconv.Atoi(getLastArg(cmd.Args()))
					if err != nil {
						return fmt.Errorf("days must be a number: %w", err)
					}
					return trash.ByDays(ctx, git, days, ProtectedBranches, TrashPrefix)
				},
			},
			{
				Name:    "empty",
				Aliases: []string{"e"},
				Usage:   "git branch -D branches with trash bin prefix",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return empty.All(ctx, git, TrashPrefix)
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func getLastArg(args cli.Args) string {
	if !args.Present() {
		return ""
	}

	return args.Get(args.Len() - 1)
}
