package rest

import (
	"github.com/bygui86/go-k8s-probes/logging"
	"github.com/bygui86/go-k8s-probes/utils"
)

const (
	restHostEnvVar = "PRODUCTS_REST_HOST"
	restPortEnvVar = "PRODUCTS_REST_PORT"

	restHostEnvVarDefault = "localhost"
	restPortEnvVarDefault = 8080
)

func loadConfig() *config {
	logging.Log.Debug("Load Products configurations")
	return &config{
		restHost: utils.GetStringEnv(restHostEnvVar, restHostEnvVarDefault),
		restPort: utils.GetIntEnv(restPortEnvVar, restPortEnvVarDefault),
	}
}
