package R

import (
	"context"
	"github.com/metrumresearchgroup/pxr/internal/configlib"
	"github.com/metrumresearchgroup/rcmd"
	"github.com/otiai10/copy"
	"os"
)

// CleanUp determines whether to clean up created directories after Failure/Success
// In the case of success, after parsing any artifacts generally the results can be cleaned up
// however sometimes may want to do additional parsing. In addition, for failures contextually
// may want to investigate, or just may want a binary yes/no in which case being
// a good citizen would mean cleaning up
type CleanUp struct {
	OnFailure bool
	OnSuccess bool
}

func NewDefaultCleanUp() CleanUp {
	return CleanUp{
		OnFailure: false,
		OnSuccess: true,
	}
}

func newTestCmdArgs(cmd string, asUser bool, interactive bool) []string {
	if cmd == "" {
		cmd = "devtools::test(stop_on_failure = TRUE)"
	}
	if interactive {
		cmd = "options(crayon.enabled = TRUE); " + cmd
	}
	cmdArgs := []string{"-e", cmd, "--slave"}
	if !asUser {
		cmdArgs = append(cmdArgs, "--vanilla")
	}
	return cmdArgs
}

// Test runs test in the testDir. If testDir is set to "", will run in
// a random tempdir.
func Test(ctx context.Context, cmd string, dir string, testDir string, cfg configlib.Config, rs rcmd.RSettings, rc rcmd.RunCfg, cleanup CleanUp) error {
	// there are situations where the users may want to run the tests in the origin dir, in which case we
	// don't want to cleanup (rm) the original directory, or the users actual files will be blown away
	if testDir == dir {
		cleanup.OnSuccess = false
		cleanup.OnFailure = false
	} else {
		// but must of the time will want to do in a temp dir
		// if the dir exists, should blow it away
		_, err := os.Lstat(testDir)
		if err != nil && !os.IsNotExist(err) {
			return err
		} else {
			// if the testDir was created by something like ioutil.TempDir() that dir would have just been created
			// to pass into this function, no need to remove and recreate
			empty, err := isEmpty(testDir)
			if !empty && err == nil {
				err = os.RemoveAll(testDir)
			}
			if err != nil {
				return err
			}
		}
		err = copy.Copy(dir, testDir)
		if err != nil {
			return err
		}
	}

	_, err := rcmd.RunR(ctx, rs, testDir, newTestCmdArgs(cmd, cfg.AsUser, true), rc)
	if err != nil {
		if cleanup.OnFailure {
			os.RemoveAll(testDir)
		}
		return err
	}
	if cleanup.OnSuccess {
		return os.RemoveAll(testDir)
	}
	return nil
}
