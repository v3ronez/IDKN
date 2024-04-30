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

func (app *application) server() error {
	serv := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", app.config.servPort),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": sig.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		err := serv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		app.logger.PrintInfo("completing background tasks", map[string]string{
			"addr": serv.Addr})

		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.PrintInfo("running server!", map[string]string{
		"mode":           app.config.envMode,
		"server_address": serv.Addr,
	})
	err := serv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError

	if err != nil {
		return err
	}
	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": serv.Addr})
	return nil
}
