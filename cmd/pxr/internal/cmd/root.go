package cmd

import (
	"context"
	"fmt"
	"github.com/metrumresearchgroup/pxr/internal/logger"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	build     string
	globalCtx context.Context
	rootCmd   = &cobra.Command{
		Use:   "pxr",
		Short: "Process eXecutoR for R",
		Long: `
interact and execute R processes and do development activities
	`,
	}
)

// Execute executes the root command.
func Execute(ctx context.Context, build string) error {
	// lets set the global context to that passed from the main cli function to everywhere else
	globalCtx = ctx
	// Execute adds all child commands to the root command sets flags appropriately.
	// This is called by main.main(). It only needs to happen once to the rootCmd.
	rootCmd.Long = fmt.Sprintf("pxr cli version %s", build)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("loglevel", "", "level for logging")
	_ = viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.SetDefault("loglevel", "info")

	rootCmd.PersistentFlags().Bool("trace", false, "use debug mode")
	_ = viper.BindPFlag("trace", rootCmd.PersistentFlags().Lookup("trace"))

}

func initConfig() {
	logger.SetLogLevel(viper.GetString("loglevel"))
	if viper.GetBool("trace") {
		viper.Debug()
		logger.SetLogLevel("trace")
	}
}
