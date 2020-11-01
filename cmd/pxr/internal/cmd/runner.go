package cmd

import (
	"context"
	"github.com/metrumresearchgroup/rcmd"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
	"os"
	"sync"
)

func run(ctx context.Context, sem *semaphore.Weighted, args []string, fn func(path string, rs rcmd.RSettings)) error {
	paths := args
	if len(paths) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal("unable to get working directory")
		}
		paths = []string{wd}
	}
	var wg sync.WaitGroup

	for _, path := range paths {
		ePath, err := homedir.Expand(path)
		if err != nil {
			log.Fatalf("error attempting to expand path: %s\n", path)
		}
		for _, lp := range cfg.LibPaths {
			if err := sem.Acquire(ctx, 1); err != nil {
				log.Errorf("error acquiring semaphore: %v", err)
				log.Errorf("Not running action for: %s", path)
				break
			}
			// semaphore will control parallelism but doesn't hold
			// the execution context from returning from the run function and exitting
			// therefore also need a waitgroup to make sure all jobs complete
			wg.Add(1)
			rs := rcmd.NewRSettings(cfg.RPath)
			for _, val := range cfg.EnvVars {
				log.Tracef("adding environment variable: %s\n", val)
				rs.EnvVars = rcmd.NvpAppendPair(rs.EnvVars, val)
			}
			rs.LibPaths = []string{lp}
			go func(rs rcmd.RSettings) {
				defer sem.Release(1)
				defer wg.Done()
				fn(ePath, rs)
			}(rs)
		}

	}
	wg.Wait()
	return nil
}
