package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	argSort              = "sort"
	argSortCommitsPushed = "commits"
	argSortWatchEvents   = "watchevents"

	argFileEvents  = "events"
	argFileCommits = "commits"
	argFileRepos   = "repos"
)

var reposCmd = &cobra.Command{
	Use:   "top-repos",
	Short: "List top 10 repos, by selected criteria.",
	Long:  "By default, sorts by the number of watch events (file provided via '--events')",
	Example: `Sort by watch events:
 gitstat top-repos --events=events.csv or
 gitstat top-repos --sort=watchevts --events=events.csv

Sort by number of pushed commits:
 gitstat top-repos --sort=commits --events=events.csv --commits=commits.csv`,
	Run: func(c *cobra.Command, args []string) {
		// parse --sort
		sort, err := c.Flags().GetString(argSort)
		cobra.CheckErr(err)

		switch sort {
		case argSortCommitsPushed,
			argSortWatchEvents:
			break
		default:
			cobra.CheckErr(errors.New("need a valid '--sort'"))
		}

		// parse --events
		events, err := c.Flags().GetString(argFileEvents)
		cobra.CheckErr(err)

		if events == "" {
			cobra.CheckErr(errors.New("need a valid '--events' path"))
		}

		// parse --repos
		repos, err := c.Flags().GetString(argFileRepos)
		cobra.CheckErr(err)

		if repos == "" {
			cobra.CheckErr(errors.New("need a valid '--repos' path"))
		}

		// maybe parse --commits
		var commits string
		if sort == argSortCommitsPushed {
			commits, err = c.Flags().GetString(argFileCommits)
			cobra.CheckErr(err)

			if commits == "" {
				cobra.CheckErr(errors.New("need a valid '--commits' path"))
			}
		}
	},
}

func init() {
	reposCmd.Flags().String(argSort, argSortWatchEvents, fmt.Sprintf("(%s | %s)", argSortWatchEvents, argSortCommitsPushed))
	reposCmd.Flags().String(argFileEvents, "", "event stream file (required)")
	reposCmd.Flags().String(argFileRepos, "", "repos stream file (required)")
	reposCmd.Flags().String(argFileCommits, "", "commit stream file (required if --sort=commits)")
}
