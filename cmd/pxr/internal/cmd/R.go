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
	if err := rcmd.StartR(globalCtx, rcmd.NewRSettings("R"), "", []string{}, *rcmd.NewRunConfig()); err != nil {
		panic(err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(RCmd)
}
