package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/config"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewServer(cfg *config.Config, router http.Handler, log log.LoggerInterface) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		IdleTimeout:  time.Duration(cfg.Server.Timeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		ErrorLog:     slog.NewLogLogger(log.Handler(), slog.LevelError),
	}

	// Create a channel to receive the error from the ListenAndServe() method
	shutdownError := make(chan error)

	// Goroutine to listen for OS interrupt and kill signals
	go func() {
		// Create a channel to receive the signals
		quit := make(chan os.Signal, 1)

		// Use signal.Notify() to register the channel to receive the specified signals
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Block until the signal is received
		s := <-quit

		log.Info("Shutting down server", "signal", s.String())

		// Create a context with a timeout of 30 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	log.Info("starting server", "addr", srv.Addr, "env", cfg.Application.Env)

	// Call the ListenAndServe() method on our http.NewServer struct
	// Only returning an error if it's not http.ErrServerClosed
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Otherwise, block until the shutdownError channel receives a value
	err = <-shutdownError
	if err != nil {
		return err
	}

	log.Info("stopped server", "addr", srv.Addr)

	return nil
}
