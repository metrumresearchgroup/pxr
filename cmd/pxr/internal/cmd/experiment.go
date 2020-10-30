package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/metrumresearchgroup/rcmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
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
	viper.Debug()
	fmt.Println("----- viper settings -------")
	prettyPrint(viper.AllSettings())
	fmt.Println("----- config settings -------")
	prettyPrint(cfg)
	wd, _ := os.Getwd()
	rs := rcmd.NewRSettings(cfg.RPath)
	rs.LibPaths = cfg.LibPaths
	rc := rcmd.NewRunConfig()
	cmdArgs := []string{"-e", ".libPaths()", "--slave"}
	if !cfg.AsUser {
		cmdArgs = append(cmdArgs, "--vanilla")
	}
	res, err := rcmd.RunR(context.Background(), rs, wd, cmdArgs, *rc)
	if err != nil {
		log.Fatal(err)
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
