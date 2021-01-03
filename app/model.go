package app

import (
	"database/sql"
	"io"
	"time"

	"github.com/openzipkin/zipkin-go/reporter"

	"github.com/bygui86/go-k8s-probes/kubernetes"
	"github.com/bygui86/go-k8s-probes/monitoring"
	"github.com/bygui86/go-k8s-probes/rest"
)

// Application implements kubernetes.Component
type Application struct {
	cfg *config

	enableMonitoring bool
	enableTracing    bool
	enableKubeProbes bool

	monitoringServer *monitoring.Server
	jaegerCloser     io.Closer
	zipkinReporter   reporter.Reporter
	dbInterface      *sql.DB
	productsServer   *rest.Server
	k8sProbesServer  *kubernetes.Server
}

type config struct {
	dbHealthCheckTimeout   time.Duration // in seconds
	restHealthCheckTimeout time.Duration // in seconds
}
