package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  1 * time.Minute,
		Handler:      app.routes(),
	}

	app.logger.Info("Server is starting", nil)

	shutdownErr := make(chan error, 1)
	go func() {
		relay := make(chan os.Signal, 1)
		signal.Notify(relay, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		s := <-relay
		app.logger.Info("Received signal:", map[string]string{
			"sig": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownErr <- err
		}

		app.logger.Info("Starting to clean up goroutines", nil)

		app.wg.Wait()
		shutdownErr <- err
	}()

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	app.logger.Info("The server has just stopped now", nil)

	err = <-shutdownErr
	if err != nil {
		app.logger.Error(err, nil)
		return err
	}

	app.logger.Info("The server has completely stopped", nil)

	return nil
}
