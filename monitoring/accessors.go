package monitoring

import (
	"github.com/gorilla/mux"
)

func (s *Server) GetRestHost() string {
	return s.config.restHost
}

func (s *Server) GetRestPort() int {
	return s.config.restPort
}

func (s *Server) GetRestRouter() *mux.Router {
	return s.router
}

func (s *Server) GetMetricsEndpoint() string {
	return metricsEndpoint
}
