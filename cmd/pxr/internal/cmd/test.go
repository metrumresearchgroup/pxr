package cmd

import (
	"fmt"
	"github.com/metrumresearchgroup/pxr/internal/R"
	"github.com/metrumresearchgroup/rcmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/semaphore"
	"io/ioutil"
	"path/filepath"
)

// checkCmd represents the R CMD check command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "run R test with the cli, in parallel",
	Long: `
	JUST FOR EXPERIMENTATION
 `,
	RunE: rTest,
}

func rTest(cmd *cobra.Command, args []string) error {
	sem := semaphore.NewWeighted(int64(cfg.Threads))
	err := run(globalCtx, sem, args, func(path string, rs rcmd.RSettings) {
		baseName := filepath.Base(rs.LibPaths[0])
		tdir, err := ioutil.TempDir("", fmt.Sprintf("*-test-%s", baseName))
		if err != nil {
			log.Fatal("error with temp dir: ", err)
		}
		log.Info(fmt.Sprintf("%s testing in directory: %s", baseName, tdir))
		if cfg.TestCmd != "" {
			log.Info(fmt.Sprintf("with command: %s", cfg.TestCmd))
		}

		cleanup := R.NewDefaultCleanUp()
		if cfg.NoCleanup {
			cleanup.OnSuccess = false
			cleanup.OnFailure = false
		}

		if err := R.Test(globalCtx,
			cfg.TestCmd,
			path,
			tdir,
			cfg,
			rs,
			*rcmd.NewRunConfig(rcmd.WithPrefix(fmt.Sprintf("[%s] ", baseName))),
			cleanup,
		); err != nil {
			log.Error(err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func init() {
	testCmd.Flags().String("test-cmd", "", "test command to run")
	viper.BindPFlag("testcmd", testCmd.Flags().Lookup("test-cmd"))
	rootCmd.AddCommand(testCmd)
}
