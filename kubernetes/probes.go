package kubernetes

import (
	"context"
	"time"

	"github.com/bygui86/go-k8s-probes/logging"
)

func New(component Component) *Server {
	logging.Log.Info("Create new Kubernetes server")

	logging.Log.Debug("Load Kubernetes configurations")
	cfg := loadConfig()

	logging.Log.Debug("Create Kubernetes server")
	kubeServer := &Server{
		config:    cfg,
		component: component,
	}
	kubeServer.setupRouter()
	kubeServer.setupHTTPServer()
	return kubeServer
}

func (s *Server) Start() {
	logging.Log.Info("Start Kubernetes server")

	if s.httpServer != nil && !s.running {
		go func() {
			err := s.httpServer.ListenAndServe()
			if err != nil {
				logging.SugaredLog.Errorf("Kubernetes server start failed: %s", err.Error())
			}
		}()
		s.running = true
		logging.SugaredLog.Infof("Kubernetes server listen on port %d", s.config.restPort)
		return
	}

	logging.Log.Error("Kubernetes server start failed: HTTP server not initialized or HTTP server already running")
}

func (s *Server) Shutdown(timeout time.Duration) {
	logging.SugaredLog.Warnf("Shutdown Kubernetes server, timeout %.0f seconds", timeout.Seconds())

	if s.httpServer != nil && s.running {
		// create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		// does not block if no connections, otherwise wait until the timeout deadline
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			logging.SugaredLog.Errorf("Kubernetes server shutdown failed: %s", err.Error())
		}
		s.running = false
		return
	}

	logging.Log.Error("Kubernetes server shutdown failed: HTTP server not initialized or HTTP server not running")
}
