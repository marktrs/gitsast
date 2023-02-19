package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/gommon/log"
	"github.com/marktrs/gitsast/internal/config"
	"github.com/marktrs/gitsast/internal/recover"
)

func Start(cfg *config.AppConfig, r http.Handler) error {
	httpLn, err := net.Listen(cfg.Server.Network, cfg.Server.Host+":"+cfg.Server.Port)
	if err != nil {
		panic(err)
	}

	handler := http.Handler(r)
	handler = recover.PanicHandler{Next: handler}

	httpServer := &http.Server{
		Addr:         httpLn.Addr().String(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		Handler:      handler,
	}

	log.Infof("listening on : %s", httpServer.Addr)

	serverErr := make(chan error, 1)
	// start server
	go func() {
		serverErr <- httpServer.Serve(httpLn)
	}()

	// set shutdown signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Errorf("server error : %v", err)
		return err
	case sig := <-shutdown:
		log.Infof("start shutdown with signal: %v", sig)
		// set deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
		defer cancel()

		// asking listener to gracefully shutdown
		err := httpServer.Shutdown(ctx)
		if err != nil {
			log.Errorf("graceful shutdown did not complete in %v : %v", cfg.Server.ShutdownTimeout.String(), err)
			return err
		}

		//  handles the response to a shutdown signal
		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return fmt.Errorf("could not stop server gracefully: %v", err)
		default:
			log.Info("shutdown complete")
		}
	}

	return nil
}
