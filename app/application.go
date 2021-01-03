package app

import (
	"io"
	"time"

	"github.com/openzipkin/zipkin-go/reporter"

	"github.com/bygui86/go-k8s-probes/commons"
	generalCfg "github.com/bygui86/go-k8s-probes/config"
	"github.com/bygui86/go-k8s-probes/logging"
	"github.com/bygui86/go-k8s-probes/monitoring"
)

func New(generalCfg *generalCfg.Config) (*Application, error) {
	logging.SugaredLog.Infof("Create new %s", commons.ServiceName)

	logging.Log.Debug("Load application configurations")
	cfg := loadConfig()

	app := &Application{
		cfg:              cfg,
		enableKubeProbes: generalCfg.GetEnableKubeProbes(),
		enableMonitoring: generalCfg.GetEnableMonitoring(),
		enableTracing:    generalCfg.GetEnableTracing(),
	}

	var monitoringServer *monitoring.Server
	if generalCfg.GetEnableMonitoring() {
		var monErr error
		monitoringServer, monErr = createMonitoring()
		if monErr != nil {
			return nil, monErr
		}
	}

	var jaegerCloser io.Closer
	var zipkinReporter reporter.Reporter
	if generalCfg.GetEnableTracing() {
		var err error
		jaegerCloser, err = initJaegerTracer()
		if err != nil {
			return nil, err
		}
	}

	dbInterface, dbErr := initDb(generalCfg.GetEnableTracing())
	if dbErr != nil {
		return nil, dbErr
	}

	prodServer, prodErr := createProducts(dbInterface)
	if prodErr != nil {
		return nil, prodErr
	}

	app.monitoringServer = monitoringServer
	app.jaegerCloser = jaegerCloser
	app.zipkinReporter = zipkinReporter
	app.dbInterface = dbInterface
	app.productsServer = prodServer

	kubeServer, kubeErr := createKubeProbes(app)
	if kubeErr != nil {
		return nil, kubeErr
	}

	app.k8sProbesServer = kubeServer

	return app, nil
}

func (a *Application) Start() error {
	logging.SugaredLog.Infof("Start %s", commons.ServiceName)

	if a.enableMonitoring {
		a.monitoringServer.Start()
	}

	err := a.productsServer.Start()
	if err != nil {
		return err
	}

	if a.enableKubeProbes {
		a.k8sProbesServer.Start()
	}

	return nil
}

func (a *Application) Shutdown(timeout time.Duration) {
	logging.SugaredLog.Warnf("Shutdown %s", commons.ServiceName)

	if a.k8sProbesServer != nil {
		a.k8sProbesServer.Shutdown(timeout)
	}

	if a.monitoringServer != nil {
		a.monitoringServer.Shutdown(timeout)
	}

	if a.jaegerCloser != nil {
		err := a.jaegerCloser.Close()
		if err != nil {
			logging.SugaredLog.Errorf("Jaeger Tracer closing failed: %s", err.Error())
		}
	}

	if a.zipkinReporter != nil {
		err := a.zipkinReporter.Close()
		if err != nil {
			logging.SugaredLog.Errorf("Zipkin Reporter closing failed: %s", err.Error())
		}
	}

	if a.productsServer != nil {
		a.productsServer.Shutdown(timeout)
	}

	if a.dbInterface != nil {
		err := a.dbInterface.Close()
		if err != nil {
			logging.SugaredLog.Errorf("DB interface closing failed: %s", err.Error())
		}
	}

	logging.SugaredLog.Warnf("Start %.0f seconds of graceful shutdown period", timeout.Seconds()+1)
	time.Sleep(timeout + 1)
}
