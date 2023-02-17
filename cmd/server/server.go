package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

func Start() error {
	serverErr := make(chan error, 1)
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Logger
	e.Logger.SetLevel(log.INFO)

	// start server
	go func() {
		serverErr <- e.Start(":1323")
	}()

	// set shutdown signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	deadline := 10 * time.Second

	select {
	case err := <-serverErr:
		e.Logger.Infof("server error : %v", err)
		return err
	case sig := <-shutdown:
		e.Logger.Infof("start shutdown with signal: %v", sig)
		// set deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), deadline)
		defer cancel()

		// asking listener to shutdown
		err := e.Shutdown(ctx)
		if err != nil {
			e.Logger.Infof("graceful shutdown did not complete in %v : %v", deadline.String(), err)
			err = e.Close()
			return err
		}

		//  handles the response to a shutdown signal
		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		default:
			e.Logger.Infof("shutdown complete")
		}
	}

	return nil
}
