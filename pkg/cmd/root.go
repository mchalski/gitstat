package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitstat",
	Short: "github event stream statistics",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(reposCmd)
	rootCmd.AddCommand(usersCmd)
}
