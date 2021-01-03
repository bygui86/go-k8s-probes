package kubernetes

const (
	// endpoints
	livenessEndpoint  = "/live"
	readinessEndpoint = "/ready"

	// response status
	ResponseStatusOk    Status = "OK"
	ResponseStatusError Status = "ERROR"

	// response codes
	ResponseCodeOk    Code = 200
	ResponseCodeError Code = 500

	// headers
	headerContentTypeKey     = "Content-Type"
	headerContentTypeAppJson = "application/json"
)
