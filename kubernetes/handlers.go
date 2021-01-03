package kubernetes

import (
	"encoding/json"
	"net/http"

	"github.com/bygui86/go-k8s-probes/logging"
)

// TODO improve error handling

func (s *Server) livenessHandler(writer http.ResponseWriter, request *http.Request) {
	logging.Log.Debug("Liveness probe invoked")

	err := json.NewEncoder(writer).Encode(s.buildProbes(writer))
	if err != nil {
		logging.SugaredLog.Errorf("JSON-Encoding liveness probe failed: %s", err.Error())
	}
}

func (s *Server) readinessHandler(writer http.ResponseWriter, request *http.Request) {
	logging.Log.Debug("Readiness probe invoked")

	err := json.NewEncoder(writer).Encode(s.buildProbes(writer))
	if err != nil {
		logging.SugaredLog.Errorf("JSON-Encoding readiness probe failed: %s", err.Error())
	}
}

func (s *Server) buildProbes(writer http.ResponseWriter) *Probe {
	logging.Log.Debug("Build Kubernetes probes")

	var probes map[string]*ComponentProbe
	if s.component != nil {
		probes = s.component.CheckStatus()
	}
	globalStatus, globalCode := computeGlobalStatus(probes)

	writer.Header().Set(headerContentTypeKey, headerContentTypeAppJson)
	probe := &Probe{
		Status:     globalStatus,
		Code:       globalCode,
		Components: probes,
	}
	return probe
}

func computeGlobalStatus(components map[string]*ComponentProbe) (Status, Code) {
	logging.Log.Debug("Compute global status")
	globalStatus := ResponseStatusOk
	globalCode := ResponseCodeOk

	if components != nil && len(components) > 0 {
		for compName, compStatus := range components {
			logging.SugaredLog.Debugf("Check %s component status", compName)
			if compStatus.IsRequired && compStatus.Code != ResponseCodeOk {
				globalStatus = compStatus.Status
				globalCode = compStatus.Code
			}
		}
	} else {
		logging.Log.Warn("No components configured to check for Kubernetes probes, returning OK per default")
	}

	return globalStatus, globalCode
}
