package config

func (c *Config) GetEnableKubeProbes() bool {
	return c.enableKubeProbes
}

func (c *Config) GetEnableMonitoring() bool {
	return c.enableMonitoring
}

func (c *Config) GetEnableTracing() bool {
	return c.enableTracing
}

func (c *Config) GetShutdownTimeout() int {
	return c.shutdownTimeout
}
