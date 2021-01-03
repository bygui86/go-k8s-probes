package config

import (
	"github.com/bygui86/go-k8s-probes/logging"
	"github.com/bygui86/go-k8s-probes/utils"
)

const (
	enableKubeProbesEnvVar = "ENABLE_KUBE_PROBES" // bool
	enableMonitoringEnvVar = "ENABLE_MONITORING"  // bool
	enableTracingEnvVar    = "ENABLE_TRACING"     // bool
	shutdownTimeoutEnvVar  = "SHUTDOWN_TIMEOUT"   // in seconds

	enableKubeProbesDefault = true
	enableMonitoringDefault = true
	enableTracingDefault    = true
	shutdownTimeoutDefault  = 10
)

func LoadConfig() *Config {
	logging.Log.Info("Load global configurations")

	shutdownTimeout := utils.GetIntEnv(shutdownTimeoutEnvVar, shutdownTimeoutDefault)
	if shutdownTimeout < 1 {
		logging.SugaredLog.Warnf("Shutdown timeout must be greater or equal to 1, fallback to default %d",
			shutdownTimeoutDefault)
		shutdownTimeout = shutdownTimeoutDefault
	}

	return &Config{
		enableKubeProbes: utils.GetBoolEnv(enableKubeProbesEnvVar, enableKubeProbesDefault),
		enableMonitoring: utils.GetBoolEnv(enableMonitoringEnvVar, enableMonitoringDefault),
		enableTracing:    utils.GetBoolEnv(enableTracingEnvVar, enableTracingDefault),
		shutdownTimeout:  shutdownTimeout,
	}
}
