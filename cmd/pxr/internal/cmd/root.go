package cmd

import (
	"context"
	"fmt"
	"github.com/metrumresearchgroup/pxr/internal/configlib"
	"github.com/metrumresearchgroup/pxr/internal/logger"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	build     string
	cfg       configlib.Config
	globalCtx context.Context
	rootCmd   = &cobra.Command{
		Use:   "pxr",
		Short: "Process eXecutor for R",
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

	rootCmd.PersistentFlags().String("loglevel", cfg.LogLevel, "level for logging")
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))

	rootCmd.PersistentFlags().StringSlice("libpath", []string{}, "path to each libpaths to run the command(s) against")
	viper.BindPFlag("libpaths", rootCmd.PersistentFlags().Lookup("libpath"))

	rootCmd.PersistentFlags().StringSlice("env", []string{}, "environment variables to set, given pattern `var=val`")
	viper.BindPFlag("environment_variables", rootCmd.PersistentFlags().Lookup("env"))

	rootCmd.PersistentFlags().String("rpath", cfg.LogLevel, "path to R")
	viper.BindPFlag("rpath", rootCmd.PersistentFlags().Lookup("rpath"))

	rootCmd.PersistentFlags().Bool("as-user", false, "use debug mode")
	viper.BindPFlag("asuser", rootCmd.PersistentFlags().Lookup("as-user"))

	rootCmd.PersistentFlags().Bool("no-cleanup", false, "use debug mode")
	viper.BindPFlag("nocleanup", rootCmd.PersistentFlags().Lookup("no-cleanup"))
}

func initConfig() {
	cfg = configlib.NewConfig(viper.GetString("config"))
	logger.SetLogLevel(cfg.LogLevel)
}
