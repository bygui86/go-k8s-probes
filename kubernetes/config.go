package kubernetes

import (
	"github.com/bygui86/go-k8s-probes/logging"
	"github.com/bygui86/go-k8s-probes/utils"
)

const (
	hostEnvVar = "KUBE_PROBES_HOST"
	portEnvVar = "KUBE_PROBES_PORT"

	hostDefault = "localhost"
	portDefault = 9091
)

func loadConfig() *config {
	logging.Log.Debug("Load Kubernetes configurations")

	return &config{
		restHost: utils.GetStringEnv(hostEnvVar, hostDefault),
		restPort: utils.GetIntEnv(portEnvVar, portDefault),
	}
}
