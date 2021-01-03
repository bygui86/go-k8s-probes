package monitoring

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	config     *config
	router     *mux.Router
	httpServer *http.Server
	running    bool
}

type config struct {
	restHost string
	restPort int
}
