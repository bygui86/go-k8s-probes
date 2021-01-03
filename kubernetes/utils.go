package kubernetes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bygui86/go-k8s-probes/commons"
	"github.com/bygui86/go-k8s-probes/logging"
)

func (s *Server) setupRouter() {
	logging.Log.Debug("Setup new Kubernetes router")

	s.router = mux.NewRouter().StrictSlash(true)
	s.router.HandleFunc(livenessEndpoint, s.livenessHandler)
	s.router.HandleFunc(readinessEndpoint, s.readinessHandler)
}

func (s *Server) setupHTTPServer() {
	logging.SugaredLog.Debugf("Setup new Kubernetes HTTP server on port %d", s.config.restPort)

	if s.config != nil {
		s.httpServer = &http.Server{
			Addr:    fmt.Sprintf(commons.HttpServerHostFormat, s.config.restHost, s.config.restPort),
			Handler: s.router,
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: commons.HttpServerWriteTimeoutDefault,
			ReadTimeout:  commons.HttpServerReadTimeoutDefault,
			IdleTimeout:  commons.HttpServerIdelTimeoutDefault,
		}
		return
	}

	logging.Log.Error("Kubernetes HTTP server creation failed: configurations not loaded")
}
