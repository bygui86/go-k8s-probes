package app

import (
	"database/sql"
	"errors"
	"io"

	"github.com/bygui86/go-k8s-probes/database"
	"github.com/bygui86/go-k8s-probes/kubernetes"
	"github.com/bygui86/go-k8s-probes/logging"
	"github.com/bygui86/go-k8s-probes/monitoring"
	"github.com/bygui86/go-k8s-probes/rest"
	"github.com/bygui86/go-k8s-probes/tracing"
)

func initDb(enableTracing bool) (*sql.DB, error) {
	logging.Log.Debug("Create new DB interface")

	var db *sql.DB
	var dbErr error
	if enableTracing {
		db, dbErr = database.NewWithWrappedTracing()
	} else {
		db, dbErr = database.New()
	}
	if dbErr != nil {
		return nil, dbErr
	}

	return db, nil
}

func createProducts(dbInterface *sql.DB) (*rest.Server, error) {
	logging.Log.Debug("Create new Products server")
	return rest.New(dbInterface)
}

func (a *Application) startProducts() error {
	logging.Log.Info("Start Products server")
	err := a.productsServer.Start()
	if err != nil {
		return err
	}

	logging.Log.Info("Products server successfully started")
	return nil
}

func createMonitoring() (*monitoring.Server, error) {
	logging.Log.Debug("Create new Monitoring server")
	server := monitoring.New()
	if server == nil {
		return nil, errors.New("monitoring server creation failed")
	}
	return server, nil
}

func (a *Application) startMonitoring() {
	logging.Log.Info("Start Monitoring server")
	a.monitoringServer.Start()
	logging.Log.Info("Monitoring server successfully started")

	rest.RegisterCustomMetrics()
	logging.Log.Info("Custom metrics successfully registered")
}

func initJaegerTracer() (io.Closer, error) {
	logging.Log.Info("Initialize Jaeger Tracer")
	closer, err := tracing.InitTracer()
	if err != nil {
		return nil, err
	}

	logging.Log.Info("Jaeger Tracer successfully initialized")
	return closer, nil
}

func createKubeProbes(app *Application) (*kubernetes.Server, error) {
	logging.Log.Debug("Create new Kubernetes server")
	server := kubernetes.New(app)
	if server == nil {
		return nil, errors.New("kubernetes server creation failed")
	}

	return server, nil
}

func (a *Application) startKubeProbes() {
	logging.Log.Info("Start Kubernetes server")
	a.k8sProbesServer.Start()
	logging.Log.Info("Kubernetes server successfully started")
}
