package app

import (
	"time"

	"github.com/bygui86/go-k8s-probes/logging"
	"github.com/bygui86/go-k8s-probes/utils"
)

const (
	dbHealthCheckTimeoutEnvVar   = "DB_HEALTH_CHECK_TIMEOUT"   // in seconds
	restHealthCheckTimeoutEnvVar = "REST_HEALTH_CHECK_TIMEOUT" // in seconds

	dbHealthCheckTimeoutDefault   = 5
	restHealthCheckTimeoutDefault = 5
)

func loadConfig() *config {
	logging.Log.Debug("Load Application configurations")
	return &config{
		dbHealthCheckTimeout: time.Duration(
			utils.GetIntEnv(dbHealthCheckTimeoutEnvVar, dbHealthCheckTimeoutDefault),
		) * time.Second,
		restHealthCheckTimeout: time.Duration(
			utils.GetIntEnv(restHealthCheckTimeoutEnvVar, restHealthCheckTimeoutDefault),
		) * time.Second,
	}
}
