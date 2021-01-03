package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bygui86/go-k8s-probes/commons"
	"github.com/bygui86/go-k8s-probes/logging"
)

const (
	// urls
	v1Endpoint         = "/api/v1"
	productsEndpoint   = v1Endpoint + "/products"
	productsIdEndpoint = productsEndpoint + "/{id:[0-9]+}"

	contentTypeHeaderKey       = "Content-Type"
	contentTypeApplicationJson = "application/json"
)

// SERVER

func (s *Server) setupRouter() {
	logging.Log.Debug("Create new router")

	s.router = mux.NewRouter().StrictSlash(true)
	s.router.HandleFunc(productsEndpoint, s.getProducts).Methods(http.MethodGet)
	s.router.HandleFunc(productsIdEndpoint, s.getProduct).Methods(http.MethodGet)
	s.router.HandleFunc(productsEndpoint, s.createProduct).Methods(http.MethodPost)
	s.router.HandleFunc(productsIdEndpoint, s.updateProduct).Methods(http.MethodPut)
	s.router.HandleFunc(productsIdEndpoint, s.deleteProduct).Methods(http.MethodDelete)
}

func (s *Server) setupHTTPServer() {
	logging.SugaredLog.Debugf("Create new HTTP server on port %d", s.config.restPort)

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

	logging.Log.Error("HTTP server creation failed: Products server configurations not loaded")
}

// HANDLERS

func sendJsonResponse(writer http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	writer.Header().Set(contentTypeHeaderKey, contentTypeApplicationJson)
	writer.WriteHeader(code)
	_, err := writer.Write(response)
	if err != nil {
		logging.SugaredLog.Errorf("Error sending JSON response: %s", err.Error())
	}
}

func sendErrorResponse(writer http.ResponseWriter, code int, message string) {
	sendJsonResponse(writer, code, map[string]string{"error": message})
}
