package kubernetes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	config     *config
	router     *mux.Router
	httpServer *http.Server
	running    bool
	component  Component
}

type config struct {
	restHost string
	restPort int
}

type Component interface {
	CheckStatus() map[string]*ComponentProbe
}

type Probe struct {
	Status     Status                     `json:"status"`
	Code       Code                       `json:"code"`
	Components map[string]*ComponentProbe `json:"components"`
}

type ComponentProbe struct {
	Status       Status  `json:"status"`
	Code         Code    `json:"code"`
	Message      string  `json:"message"`
	TimeConsumed float64 `json:"timeConsumed"` // in milliseconds
	IsRequired   bool    `json:"isRequired"`
}

type Status string
type Code int
