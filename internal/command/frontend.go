package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/finance-planner/internal/errs"
	"github.com/tomwright/finance-planner/internal/frontend"
	"github.com/tomwright/finance-planner/internal/util/shutdownutil"
	"sync"
)

func HTTPFrontend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "frontend",
		Short: "Run a HTTP server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			listenAddress, _ := cmd.Flags().GetString("listen-address")
			baseURL, _ := cmd.Flags().GetString("base-url")
			assetsPath, _ := cmd.Flags().GetString("assets")

			// wg contains a counter for all services started in this command.
			wg := &sync.WaitGroup{}
			// If any error is written to errCh, all services will be shutdown
			// and the command will finish.
			errCh := make(chan error)
			// When shutdownCh is closed, all services will begin a graceful
			// shutdown.
			shutdownCh := make(chan struct{})

			go shutdownutil.HandleShutdownSignal(errCh)

			// Start HTTP service.
			wg.Add(1)
			go frontend.Start(assetsPath, baseURL, listenAddress, wg, errCh, shutdownCh)

			// Block until errCh message
			err := <-errCh
			fmt.Printf("Error received: %s\n", err)

			// Notify all services of shutdown
			close(shutdownCh)

			// Wait for all services to stop
			wg.Wait()

			// Only return an error from the command if a
			// non-shutdown-signal error was written to errCh.
			if err == nil {
				return nil
			}
			e := errs.FromErr(err)
			if e.Code() == errs.ErrShutdownSignal {
				return nil
			}
			return e
		},
	}

	cmd.Flags().String("listen-address", ":80", "HTTP listen address")
	cmd.Flags().String("base-url", "http://localhost:8080", "Base URL for the API")
	cmd.Flags().String("assets", "", "The path to the frontend server assets")

	_ = cmd.MarkFlagRequired("assets")

	return cmd
}
