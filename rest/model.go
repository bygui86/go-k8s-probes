package rest

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	config      *config
	router      *mux.Router
	httpServer  *http.Server
	dbInterface *sql.DB
	running     bool
}

type config struct {
	restHost string
	restPort int
}
