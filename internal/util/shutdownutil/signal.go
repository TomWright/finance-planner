package shutdownutil

import (
	"github.com/tomwright/finance-planner/internal/errs"
	"os"
	"os/signal"
	"syscall"
)

// HandleShutdownSignal will write a ErrShutdownSignal to the given errCh
// when a shutdown signal is received.
// It is expected that this func is run in a go routine.
func HandleShutdownSignal(errCh chan error) {
	quitCh := make(chan os.Signal)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGTERM)

	hit := false
	for {
		<-quitCh
		if hit {
			os.Exit(0)
		}
		if !hit {
			errCh <- errs.New().
				WithCode(errs.ErrShutdownSignal).
				WithMessage("shutdown signal received")
		}
		hit = true
	}
}
