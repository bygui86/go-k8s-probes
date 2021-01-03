package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"

	"github.com/bygui86/go-k8s-probes/commons"
	"github.com/bygui86/go-k8s-probes/database"
	"github.com/bygui86/go-k8s-probes/kubernetes"
	"github.com/bygui86/go-k8s-probes/logging"
	"github.com/bygui86/go-k8s-probes/time_measure"
)

const (
	headerAccept          = "Accept"
	headerContentType     = "Content-Type"
	headerApplicationJson = "application/json"
)

var restClient *http.Client

func (a *Application) CheckStatus() map[string]*kubernetes.ComponentProbe {
	logging.SugaredLog.Debugf("Check %s status", commons.ServiceName)

	if restClient == nil {
		restClient = &http.Client{
			Timeout: a.cfg.restHealthCheckTimeout,
		}
	}

	components := make(map[string]*kubernetes.ComponentProbe, 4)
	// required
	components["db"] = a.checkDbStatus()
	components["products"] = a.checkProductsStatus()
	components["monitoring"] = a.checkMonitoringStatus()
	// not required
	components["tracing"] = a.checkTracingStatus()

	return components
}

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
			msg = fmt.Sprintf("DB interface NOT HEALTHY: %s", err.Error())
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
		msg = fmt.Sprintf("Products API NOT HEALTHY: %s", dialErr.Error())
	} else {
		logging.Log.Debug("Walk through endpoints")
		walkErr := a.productsServer.GetRestRouter().Walk(routeWalker)
		if walkErr != nil {
			status = kubernetes.ResponseStatusError
			code = kubernetes.ResponseCodeError
			msg = fmt.Sprintf("Products API NOT HEALTHY: %s", walkErr.Error())
		} else {
			status = kubernetes.ResponseStatusOk
			code = kubernetes.ResponseCodeOk
			msg = "Products API healthy"
		}
	}

	logging.Log.Debug("Check Rest response")
	metricsErr := responseChecker(
		"Products",
		a.productsServer.GetRestHost(),
		a.productsServer.GetRestPort(),
		a.productsServer.GetProductsEndpoint(),
		map[string]string{
			headerAccept:      headerApplicationJson,
			headerContentType: headerApplicationJson,
		},
		checkProductsResponse,
	)
	if metricsErr != nil {
		status = kubernetes.ResponseStatusError
		code = kubernetes.ResponseCodeError
		msg = fmt.Sprintf("Monitoring NOT HEALTHY: %s", metricsErr.Error())
	} else {
		status = kubernetes.ResponseStatusOk
		code = kubernetes.ResponseCodeOk
		msg = "Monitoring healthy"
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

func checkProductsResponse(response *http.Response) error {
	var products []*database.Product
	unmarshErr := json.NewDecoder(response.Body).Decode(&products)
	if unmarshErr != nil {
		return unmarshErr
	}
	return nil
}

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
		msg = fmt.Sprintf("Monitoring API NOT HEALTHY: %s", dialErr.Error())
	} else {
		logging.Log.Debug("Walk through endpoints")
		walkErr := a.productsServer.GetRestRouter().Walk(routeWalker)
		if walkErr != nil {
			status = kubernetes.ResponseStatusError
			code = kubernetes.ResponseCodeError
			msg = fmt.Sprintf("Monitoring NOT HEALTHY: %s", walkErr.Error())
		} else {
			status = kubernetes.ResponseStatusOk
			code = kubernetes.ResponseCodeOk
			msg = "Monitoring healthy"
		}
	}

	logging.Log.Debug("Check Rest response")
	metricsErr := responseChecker(
		"Monitoring",
		a.monitoringServer.GetRestHost(),
		a.monitoringServer.GetRestPort(),
		a.monitoringServer.GetMetricsEndpoint(),
		map[string]string{
			// headerAccept:      headerApplicationJson,
			// headerContentType: headerApplicationJson,
		},
		checkMonitoringResponse,
	)
	if metricsErr != nil {
		status = kubernetes.ResponseStatusError
		code = kubernetes.ResponseCodeError
		msg = fmt.Sprintf("Monitoring NOT HEALTHY: %s", metricsErr.Error())
	} else {
		status = kubernetes.ResponseStatusOk
		code = kubernetes.ResponseCodeOk
		msg = "Monitoring healthy"
	}

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

func checkMonitoringResponse(response *http.Response) error {
	bodyBytes, bodyErr := ioutil.ReadAll(response.Body)
	if bodyErr != nil {
		return fmt.Errorf("Monitoring response body reading failed: %s", bodyErr.Error())
	}
	bodyString := string(bodyBytes)
	logging.SugaredLog.Debugf("Monitoring response: %s", bodyString)

	if bodyString == "" {
		return errors.New("empty Monitoring response")
	}
	return nil
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

func responseChecker(
	system string,
	host string, port int, endpoint string, headers map[string]string,
	checkResponse func(*http.Response) error) error {

	// prepare url
	baseURL, urlErr := url.Parse(fmt.Sprintf("http://%s:%d", host, port))
	if urlErr != nil {
		return urlErr
	}
	endpointUrl := &url.URL{Path: endpoint}
	path := baseURL.ResolveReference(endpointUrl)

	// prepare request
	restRequest, reqErr := http.NewRequest(http.MethodGet, path.String(), nil)
	if reqErr != nil {
		return reqErr
	}
	for key, val := range headers {
		restRequest.Header.Set(key, val)
	}

	// get response
	response, respErr := restClient.Do(restRequest)
	if respErr != nil {
		return respErr
	}
	if response.StatusCode != http.StatusOK {
		logging.SugaredLog.Debugf("%s response code %d", system, response.StatusCode)
		return fmt.Errorf("%s response code %d", system, response.StatusCode)
	}
	defer closeResponse(response)

	// check response
	return checkResponse(response)
}

func closeResponse(response *http.Response) {
	err := response.Body.Close()
	if err != nil {
		logging.SugaredLog.Warnf("Error closing rest response body: %s", err.Error())
	}
}
