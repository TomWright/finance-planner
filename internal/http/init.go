package http

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/tomwright/finance-planner/internal/application/service"
	"net"
	"net/http"
	"sync"
	"time"
)

// Start starts up a HTTP server.
// It is expected that Start will be executed in a go routine.
// wg.Add(1) should have been called already.
// If shutdownCh is closed, the server should be shutdown.
func Start(profileService service.Profile, listenAddress string, wg *sync.WaitGroup, errCh chan error, shutdownCh chan struct{}) {
	// Ensure the wg.Done() is decremented.
	defer wg.Done()

	r := chi.NewRouter()

	for _, h := range loadHandlers(profileService) {
		h.Bind(r)
	}

	server := &http.Server{
		Handler: r,
	}

	startErrCh := make(chan error)

	startFn := func() {
		listener, err := net.Listen("tcp", listenAddress)
		if err != nil {
			startErrCh <- fmt.Errorf("could not listen on address `%s`: %s", listenAddress, err)
			return
		}

		fmt.Printf("http server listening on `%s`\n", listenAddress)

		_ = server.Serve(listener)
	}
	stopFn := func() {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
		fmt.Printf("http server shutting down\n")
		err := server.Shutdown(ctx)
		if err != nil {
			errCh <- err
		}
		fmt.Printf("http server shut down\n")
	}

	// Start the HTTP server in a routine
	go startFn()

	// Wait for either a start-up error, or for the shutdownCh to be closed.
	select {
	case err := <-startErrCh:
		// No need to stopFn() because the server wasn't started up.
		errCh <- err
	case <-shutdownCh:
		// Call stopFn() to gracefully shutdown.
		stopFn()
	}
}

// loadHandlers returns all of the handlers to be served via HTTP.
func loadHandlers(profileService service.Profile) []Handler {
	return []Handler{
		NewListTransactionsHandler(profileService),
		NewAddTransactionHandler(profileService),
		NewStatsTransactionsHandler(profileService),
	}
}
