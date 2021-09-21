package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/mchalski/gitstat/pkg/actions"
	"github.com/mchalski/gitstat/pkg/data"
)

const (
	argFileActors = "actors"
)

var usersCmd = &cobra.Command{
	Use:     "top-users",
	Short:   "List top 10 users, by PRs created and commits pushed.",
	Example: "gitstat top-users",
	Run: func(c *cobra.Command, args []string) {
		// parse --events
		events, err := c.Flags().GetString(argFileEvents)
		cobra.CheckErr(err)

		if events == "" {
			cobra.CheckErr(errors.New("need a valid '--events' path"))
		}

		// parse --commits
		commits, err := c.Flags().GetString(argFileCommits)
		cobra.CheckErr(err)

		if commits == "" {
			cobra.CheckErr(errors.New("need a valid '--commits' path"))
		}

		// parse --actors
		actors, err := c.Flags().GetString(argFileActors)
		cobra.CheckErr(err)

		if actors == "" {
			cobra.CheckErr(errors.New("need a valid '--actors' path"))
		}

		streamEvts, err := data.NewCsvStream(events)
		cobra.CheckErr(err)

		streamCommits, err := data.NewCsvStream(commits)
		cobra.CheckErr(err)

		streamActors, err := data.NewCsvStream(actors)
		cobra.CheckErr(err)

		action := actions.NewTopUsers(streamEvts,
			streamCommits,
			streamActors)

		action.Run()
	},
}

func init() {
	usersCmd.Flags().String(argFileEvents, "", "event stream file (required)")
	usersCmd.Flags().String(argFileCommits, "", "commit stream file (required)")
	usersCmd.Flags().String(argFileActors, "", "actors stream file (required)")
}
