package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkCmd represents the R CMD check command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "debug with the cli",
	Long: `
	JUST FOR EXPERIMENTATION
 `,
	RunE: rDebug,
}

func rDebug(cmd *cobra.Command, args []string) error {
	viper.Debug()
	fmt.Println("----- viper settings -------")
	prettyPrint(viper.AllSettings())
	fmt.Println("----- config settings -------")
	prettyPrint(cfg)
	return nil
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
