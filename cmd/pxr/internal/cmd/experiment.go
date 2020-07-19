package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/metrumresearchgroup/rcmd"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// checkCmd represents the R CMD check command
var experimentCmd = &cobra.Command{
	Use:   "experiment",
	Short: "experiment with the cli",
	Long: `
	JUST FOR EXPERIMENTATION
 `,
	RunE: rExperiment,
}

func rExperiment(cmd *cobra.Command, args []string) error {

	dir, _ := homedir.Expand("~/metrum/metrumresearchgroup/matrixbuilds")
	res, err := rcmd.RunR(context.Background(), rcmd.NewRSettings("R"), dir, []string{"-e", "options(crayon.enabled = TRUE); devtools::test()", "--slave", "--interactive"}, *rcmd.NewRunConfig())
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	return nil
}

func init() {
	rootCmd.AddCommand(experimentCmd)
}

func prettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
