package main

import (
	"context"
	"fmt"
	"github.com/metrumresearchgroup/pxr/cmd/pxr/internal/cmd"
	"github.com/pkg/errors"
	"log"
	_ "net/http/pprof" // Register the pprof handlers
	"os"
	"os/signal"
	"syscall"
	"time"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "REPLACEME"

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {

	// =========================================================================
	// Logging

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	cliErrors := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	//
	//// run the cli
	go func() {
		cliErrors <- cmd.Execute(ctx, build)
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-cliErrors:
		return errors.Wrap(err, "cli error")

	case sig := <-shutdown:
		// need to figure out how to shutdown gracefully

		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
			//case err != nil:
			//	return errors.Wrap(err, "could not stop server gracefully")
		default:
			fmt.Println("calling cancel...")
			cancel()
			time.Sleep(2 * time.Second)
		}
	}

	return nil
}
