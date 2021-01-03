package monitoring

import (
	"github.com/bygui86/go-k8s-probes/logging"
	"github.com/bygui86/go-k8s-probes/utils"
)

const (
	monitorHostEnvVar = "MONITOR_HOST"
	monitorPortEnvVar = "MONITOR_PORT"

	monitorHostDefault = "localhost"
	monitorPortDefault = 9090
)

func loadConfig() *config {
	logging.Log.Debug("Load monitoring configurations")
	return &config{
		restHost: utils.GetStringEnv(monitorHostEnvVar, monitorHostDefault),
		restPort: utils.GetIntEnv(monitorPortEnvVar, monitorPortDefault),
	}
}
