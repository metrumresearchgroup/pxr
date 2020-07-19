package cmd

import (
	"fmt"
	"github.com/metrumresearchgroup/pxr/internal/R"
	"github.com/metrumresearchgroup/rcmd"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"sync"
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

	dir, _ := homedir.Expand("~/metrum/metrumresearchgroup/matrixbuilds")
	rs := rcmd.NewRSettings("R")
	var wg sync.WaitGroup
	wg.Add(2)

	// TODO: these should return whether they error'd to a channel to resurface at end of cmd
	go func(rs rcmd.RSettings) {
		rs.EnvVars = rcmd.NvpAppendPair(rs.EnvVars, "R_LIBS_SITE=~/rpkgs/2020-03-24")
		tdir, err := ioutil.TempDir("", "*-test-2020-03-24")
		if err != nil {
			log.Error("error with temp dir: ", err)
		}
		log.Info("[2020-03-24] testing in directory: ", tdir)
		if err := R.Test(globalCtx, dir, tdir, rs, *rcmd.NewRunConfig(rcmd.WithPrefix("[2020-03-24]")), R.NewDefaultCleanUp()); err != nil {
			log.Error(err)
		}
		wg.Done()
	}(rs)

	go func(rs rcmd.RSettings) {
		rs.EnvVars = rcmd.NvpAppendPair(rs.EnvVars, "R_LIBS_SITE=~/rpkgs/2020-06-08")
		tdir, err := ioutil.TempDir("", "*-test-2020-06-08")
		if err != nil {
			log.Fatal("error with temp dir: ", err)
		}
		log.Info("[2020-06-08] testing in directory: ", tdir)
		if err := R.Test(globalCtx, dir, tdir, rs, *rcmd.NewRunConfig(rcmd.WithPrefix("[2020-06-08]")), R.NewDefaultCleanUp()); err != nil {
			log.Error(err)
		}
		wg.Done()
	}(rs)
	wg.Wait()
	fmt.Println(prettyPrint(rs))
	return nil
}

func init() {
	rootCmd.AddCommand(testCmd)
}
