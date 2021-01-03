package commons

import "time"

const (
	ServiceName = "product-service"

	HttpServerHostFormat          = "%s:%d"
	HttpServerWriteTimeoutDefault = time.Second * 15
	HttpServerReadTimeoutDefault  = time.Second * 15
	HttpServerIdelTimeoutDefault  = time.Second * 60

	SecondDivider   = 1000000000
	MillisecDivider = 1000000
)
