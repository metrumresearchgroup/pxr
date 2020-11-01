package cmd

import (
	"fmt"
	"github.com/metrumresearchgroup/pxr/internal/R"
	"github.com/metrumresearchgroup/rcmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/semaphore"
	"io/ioutil"
	"path/filepath"
)

// checkCmd represents the R CMD check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "run R CMD check with the cli, in parallel",
	Long: `
    run R CMD check on package(s) in parallel.	
 `,
	RunE: rCheck,
}

func rCheck(cmd *cobra.Command, args []string) error {
	sem := semaphore.NewWeighted(int64(cfg.Threads))
	run(globalCtx, sem, args, func(path string, rs rcmd.RSettings) {
		rs.EnvVars = rcmd.NvpAppend(rs.EnvVars, "BABYLON_EXE_PATH", "bbi")
		baseName := filepath.Base(rs.LibPaths[0])
		tdir, err := ioutil.TempDir("", fmt.Sprintf("*-check-%s", baseName))
		if err != nil {
			log.Fatal("error with temp dir: ", err)
		}
		log.Info(fmt.Sprintf("%s checking in directory: %s", baseName, tdir))
		cleanup := R.NewDefaultCleanUp()
		if cfg.NoCleanup {
			cleanup.OnSuccess = false
			cleanup.OnFailure = false
		}
		if err := R.Check(globalCtx,
			path,
			tdir,
			cfg,
			rcmd.CheckArgs{Output: tdir},
			rs,
			*rcmd.NewRunConfig(rcmd.WithPrefix(fmt.Sprintf("[%s] ", baseName))),
			cleanup,
		); err != nil {
			log.Error(err)
		}
	})
	return nil
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
