package cmd

import (
	"github.com/metrumresearchgroup/rcmd"
	"github.com/spf13/cobra"
)

// checkCmd represents the R CMD check command
var RCmd = &cobra.Command{
	Use:   "R",
	Short: "R with the cli",
	Long: `
   Start R
 `,
	RunE: rR,
}

func rR(cmd *cobra.Command, args []string) error {
	rs := rcmd.NewRSettings("R")
	cmdArgs := []string{"--vanilla"}
	if cfg.AsUser {
		rs.AsUser = true
		// don't run vanilla so will act like user session
		cmdArgs = []string{}
	}
	if err := rcmd.StartR(globalCtx, rs, "", cmdArgs, *rcmd.NewRunConfig()); err != nil {
		panic(err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(RCmd)
}
