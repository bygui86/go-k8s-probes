package app

import (
	"context"
	"fmt"
	"net"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"

	"github.com/bygui86/go-k8s-probes/commons"
	"github.com/bygui86/go-k8s-probes/kubernetes"
	"github.com/bygui86/go-k8s-probes/logging"
	"github.com/bygui86/go-k8s-probes/time_measure"
)

func (a *Application) CheckStatus() map[string]*kubernetes.ComponentProbe {
	logging.SugaredLog.Debugf("Check %s status", commons.ServiceName)

	components := make(map[string]*kubernetes.ComponentProbe, 1)
	// required
	components["db"] = a.checkDbStatus()
	components["products"] = a.checkProductsStatus()
	components["monitoring"] = a.checkMonitoringStatus()
	// not required
	components["tracing"] = a.checkTracingStatus()

	return components
}

// TODO TBD: check response / check code only / check response and code
func (a *Application) checkDbStatus() *kubernetes.ComponentProbe {
	timeMeasure := time_measure.StartTimeMeasure()

	logging.Log.Debug("Check DB interface status")
	var status kubernetes.Status
	var code kubernetes.Code
	var msg string

	if a.dbInterface != nil {
		ctx, cancel := context.WithTimeout(context.Background(), a.cfg.dbHealthCheckTimeout)
		defer cancel()

		err := a.dbInterface.PingContext(ctx)

		if err != nil {
			status = kubernetes.ResponseStatusError
			code = kubernetes.ResponseCodeError
			msg = err.Error()
		} else {
			status = kubernetes.ResponseStatusOk
			code = kubernetes.ResponseCodeOk
			msg = "DB interface healthy"
		}
	} else {
		status = kubernetes.ResponseStatusError
		code = kubernetes.ResponseCodeError
		msg = "DB interface not initialized"
	}

	timeMeasure.StopTimeMeasure()
	milliSec, _ := timeMeasure.GetDeltaInMil().Float64()

	logging.Log.Debug("DB interface status checked")
	return &kubernetes.ComponentProbe{
		Status:       status,
		Code:         code,
		Message:      msg,
		TimeConsumed: milliSec,
		IsRequired:   true,
	}
}

func (a *Application) checkProductsStatus() *kubernetes.ComponentProbe {
	timeMeasure := time_measure.StartTimeMeasure()

	logging.Log.Debug("Check Products status")
	var status kubernetes.Status
	var code kubernetes.Code
	var msg string

	logging.Log.Debug("Build Products address")
	address := fmt.Sprintf("localhost:%d", a.productsServer.GetRestPort())

	logging.Log.Debug("Check TCP connection")
	_, dialErr := net.DialTimeout("tcp", address, a.cfg.restHealthCheckTimeout)
	if dialErr != nil {
		status = kubernetes.ResponseStatusError
		code = kubernetes.ResponseCodeError
		msg = dialErr.Error()
	} else {
		logging.Log.Debug("Walk through endpoints")
		walkErr := a.productsServer.GetRouter().Walk(routeWalker)
		if walkErr != nil {
			status = kubernetes.ResponseStatusError
			code = kubernetes.ResponseCodeError
			msg = walkErr.Error()
		} else {
			status = kubernetes.ResponseStatusOk
			code = kubernetes.ResponseCodeOk
			msg = "Products healthy"
		}
	}

	timeMeasure.StopTimeMeasure()
	milliSec, _ := timeMeasure.GetDeltaInMil().Float64()

	logging.Log.Debug("Products status checked")
	return &kubernetes.ComponentProbe{
		Status:       status,
		Code:         code,
		Message:      msg,
		TimeConsumed: milliSec,
		IsRequired:   true,
	}
}

// TODO check response and code
func (a *Application) checkMonitoringStatus() *kubernetes.ComponentProbe {
	timeMeasure := time_measure.StartTimeMeasure()

	logging.Log.Debug("Check Monitoring status")
	var status kubernetes.Status
	var code kubernetes.Code
	var msg string

	logging.Log.Debug("Build Monitoring address")
	address := fmt.Sprintf("localhost:%d", a.monitoringServer.GetRestPort())

	logging.Log.Debug("Check TCP connection")
	_, dialErr := net.DialTimeout("tcp", address, a.cfg.restHealthCheckTimeout)
	if dialErr != nil {
		status = kubernetes.ResponseStatusError
		code = kubernetes.ResponseCodeError
		msg = dialErr.Error()
	} else {
		logging.Log.Debug("Walk through endpoints")
		walkErr := a.productsServer.GetRouter().Walk(routeWalker)
		if walkErr != nil {
			status = kubernetes.ResponseStatusError
			code = kubernetes.ResponseCodeError
			msg = walkErr.Error()
		} else {
			status = kubernetes.ResponseStatusOk
			code = kubernetes.ResponseCodeOk
			msg = "Monitoring healthy"
		}
	}

	// TODO check for response!

	timeMeasure.StopTimeMeasure()
	milliSec, _ := timeMeasure.GetDeltaInMil().Float64()

	logging.Log.Debug("Monitoring status checked")
	return &kubernetes.ComponentProbe{
		Status:       status,
		Code:         code,
		Message:      msg,
		TimeConsumed: milliSec,
		IsRequired:   true,
	}
}

func (a *Application) checkTracingStatus() *kubernetes.ComponentProbe {
	timeMeasure := time_measure.StartTimeMeasure()

	logging.Log.Debug("Check Tracing status")
	var status kubernetes.Status
	var code kubernetes.Code
	var msg string

	if a.jaegerCloser != nil {
		if opentracing.IsGlobalTracerRegistered() {
			status = kubernetes.ResponseStatusOk
			code = kubernetes.ResponseCodeOk
			msg = "Jaeger Tracer registered"
		} else {
			status = kubernetes.ResponseStatusError
			code = kubernetes.ResponseCodeError
			msg = "Jaeger Tracer NOT REGISTERED"
		}
	}

	timeMeasure.StopTimeMeasure()
	milliSec, _ := timeMeasure.GetDeltaInMil().Float64()

	logging.Log.Debug("Tracing status checked")
	return &kubernetes.ComponentProbe{
		Status:       status,
		Code:         code,
		Message:      msg,
		TimeConsumed: milliSec,
		IsRequired:   false,
	}
}

func routeWalker(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	logging.SugaredLog.Debugf("Check route: %s", route.GetName())
	if route.GetHandler() == nil {
		return fmt.Errorf("route %s handler is nil", route.GetName())
	}
	return nil
}

// TODO
// func responseChecker() error {
// 	// prepare request
// 	endpointUrl := &url.URL{Path: rootProductsEndpoint}
// 	path := s.baseURL.ResolveReference(endpointUrl)
// 	restRequest, reqErr := http.NewRequest(http.MethodGet, path.String(), nil)
// 	if reqErr != nil {
// 		return reqErr
// 	}
// 	restRequest.Header.Set(headerAccept, headerApplicationJson)
// 	restRequest.Header.Set(headerUserAgent, headerUserAgentClient)
//
// 	// get response
// 	response, respErr := s.restClient.Do(restRequest)
// 	if respErr != nil {
// 		return respErr
// 	}
//
// }
