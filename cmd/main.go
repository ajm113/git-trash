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
		Usage: "Moves branches to a trash bin by moving targeted branches for trash by adding a prefix (" + TrashPrefix + "). Perfect for large git projects that handle many branches",
		Commands: []*cli.Command{
			{
				Name:    "trash",
				Aliases: []string{"t"},
				Usage:   "Moves targeted branch(es) to trash bin. This also accepts glob characters such as * and ?.",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return trash.ByMatch(ctx, git, getLastArg(cmd.Args()), ProtectedBranches, TrashPrefix)
				},
			},
			{
				Name:    "days",
				Aliases: []string{"d"},
				Usage:   "Moves branches older then the number of days provided to the trash bin.",
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
				Usage:   "Emptys the trash bin",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return empty.All(ctx, git, TrashPrefix)
				},
			},
			{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "Shows licensing, download, and issue reporting information.",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Printf("git-trash %s (commit %s, built %s)\n\n", Version, Commit, Date)
					fmt.Println("License:")
					fmt.Println("  Licensed under the GNU General Public License v3.0 (GPL-3.0).")
					fmt.Println("  This program comes with ABSOLUTELY NO WARRANTY. It is free software,")
					fmt.Println("  and you are welcome to redistribute it under certain conditions.")
					fmt.Println("  See https://www.gnu.org/licenses/gpl-3.0.html for the full text.")
					fmt.Println()
					fmt.Println("Download / Source:")
					fmt.Println("  https://github.com/ajm113/git-trash")
					fmt.Println()
					fmt.Println("Report Issues:")
					fmt.Println("  https://github.com/ajm113/git-trash/issues")
					return nil
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
