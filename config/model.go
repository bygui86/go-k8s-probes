package config

type Config struct {
	enableKubeProbes bool
	enableMonitoring bool
	enableTracing    bool
	shutdownTimeout  int
}
