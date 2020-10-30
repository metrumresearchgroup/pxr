package R

import (
	"context"
	"github.com/metrumresearchgroup/pxr/internal/configlib"
	"github.com/metrumresearchgroup/rcmd"
	"github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// Check runs test in the testDir. If testDir is set to "", will run in
// a random tempdir.
//   Check R packages from package sources, which can be directories or
//   package 'tar' archives with extension '.tar.gz', '.tar.bz2',
//   '.tar.xz' or '.tgz'
func Check(ctx context.Context, path string, testDir string, cfg configlib.Config, cs rcmd.CheckArgs, rs rcmd.RSettings, rc rcmd.RunCfg, cleanup CleanUp) error {
	// there are situations where the users may want to run the tests in the origin dir, in which case we
	// don't want to cleanup (rm) the original directory, or the users actual files will be blown away
	if testDir == path {
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
		log.Debug(rs)
		err = copy.Copy(path, filepath.Join(testDir, filepath.Base(path)))
		if err != nil {
			return err
		}
	}

	var cmdArgs []string
	if !cfg.AsUser {
		cmdArgs = append(cmdArgs, "--vanilla")
	}
	cmdArgs = append(
		cmdArgs,
		"CMD",
		"check",
		filepath.Join(testDir, filepath.Base(path)),
	)
	cmdArgs = append(cmdArgs, cs.CliArgs()...)
	log.Debug("running R with command: %s", cmdArgs)
	_, err := rcmd.RunR(ctx, rs, testDir, cmdArgs, rc)
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
