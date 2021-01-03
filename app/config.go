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

	dbTimeout := utils.GetIntEnv(dbHealthCheckTimeoutEnvVar, dbHealthCheckTimeoutDefault)
	if dbTimeout > 0 {
		logging.SugaredLog.Warnf("DB health check timeout must be greater than 0, fallback to default %d",
			dbHealthCheckTimeoutDefault)
		dbTimeout = dbHealthCheckTimeoutDefault
	}

	restTimeout := utils.GetIntEnv(restHealthCheckTimeoutEnvVar, restHealthCheckTimeoutDefault)
	if restTimeout > 0 {
		logging.SugaredLog.Warnf("Rest health check timeout must be greater than 0, fallback to default %d",
			restHealthCheckTimeoutDefault)
		restTimeout = restHealthCheckTimeoutDefault
	}

	return &config{
		dbHealthCheckTimeout:   time.Duration(dbTimeout) * time.Second,
		restHealthCheckTimeout: time.Duration(restTimeout) * time.Second,
	}
}
