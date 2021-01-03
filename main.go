package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bygui86/go-k8s-probes/app"
	"github.com/bygui86/go-k8s-probes/commons"
	"github.com/bygui86/go-k8s-probes/config"
	"github.com/bygui86/go-k8s-probes/logging"
)

func main() {
	initLogging()

	logging.SugaredLog.Infof("Load %s configurations", commons.ServiceName)
	cfg := loadConfig()

	logging.SugaredLog.Infof("Create %s", commons.ServiceName)
	application, newErr := app.New(cfg)
	if newErr != nil {
		logging.SugaredLog.Errorf("%s creationg failed: %s", commons.ServiceName, newErr.Error())
		os.Exit(502)
	}

	logging.SugaredLog.Infof("Start %s", commons.ServiceName)
	startErr := application.Start()
	if startErr != nil {
		logging.SugaredLog.Errorf("%s start failed: %s", commons.ServiceName, startErr.Error())
		os.Exit(503)
	}
	logging.SugaredLog.Infof("%s up and running", commons.ServiceName)

	startSysCallChannel()

	logging.SugaredLog.Infof("Shutdown %s", commons.ServiceName)
	application.Shutdown(time.Duration(cfg.GetShutdownTimeout()) * time.Second)
}

func initLogging() {
	err := logging.InitGlobalLogger()
	if err != nil {
		logging.SugaredLog.Errorf("Logging initialization failed: %s", err.Error())
		os.Exit(501)
	}
}

func loadConfig() *config.Config {
	logging.Log.Debug("Load configurations")
	return config.LoadConfig()
}

func startSysCallChannel() {
	syscallCh := make(chan os.Signal)
	signal.Notify(syscallCh, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-syscallCh
}
