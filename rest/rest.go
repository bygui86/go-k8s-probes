package rest

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bygui86/go-k8s-probes/logging"
)

func New(dbInterface *sql.DB) (*Server, error) {
	logging.Log.Info("Create new Products server")

	cfg := loadConfig()

	server := &Server{
		config:      cfg,
		dbInterface: dbInterface,
	}

	server.setupRouter()
	server.setupHTTPServer()
	return server, nil
}

func (s *Server) Start() error {
	logging.Log.Info("Start Products server")

	if s.httpServer != nil && !s.running {
		var err error
		go func() {
			err = s.httpServer.ListenAndServe()
			if err != nil {
				logging.SugaredLog.Errorf("Products server start failed: %s", err.Error())
			}
		}()
		if err != nil {
			return err
		}
		s.running = true
		logging.SugaredLog.Infof("Products server listening on port %d", s.config.restPort)
		return nil
	}

	return fmt.Errorf("products server start failed: HTTP server not initialized or HTTP server already running")
}

func (s *Server) Shutdown(timeout time.Duration) {
	logging.SugaredLog.Warnf("Shutdown Products server, timeout %.0f seconds", timeout.Seconds())

	if s.httpServer != nil && s.running {
		// create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		// does not block if no connections, otherwise wait until the timeout deadline
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			logging.SugaredLog.Errorf("Products server shutdown failed: %s", err.Error())
		}

		s.running = false
		return
	}

	logging.Log.Error("Products server shutdown failed: HTTP server not initialized or HTTP server not running")
}
