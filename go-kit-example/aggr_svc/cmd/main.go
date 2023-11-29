package main

import (
	"net"
	"net/http"
	"os"

	"github.com/go-kit/log"
	aggrendpoint "github.com/raphaelmb/go-toll-calculator/go-kit-example/aggr_svc/aggr_endpoint"
	aggrservice "github.com/raphaelmb/go-toll-calculator/go-kit-example/aggr_svc/aggr_service"
	aggrtransport "github.com/raphaelmb/go-toll-calculator/go-kit-example/aggr_svc/aggr_transport"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	service := aggrservice.New()
	endpoints := aggrendpoint.New(service, logger)
	httpHandler := aggrtransport.NewHTTPHandler(endpoints, logger)

	httpListener, err := net.Listen("tcp", ":3000")
	if err != nil {
		logger.Log("transport", "HTTP", "during", "Listen", "err", err)
		os.Exit(1)
	}
	logger.Log("transport", "HTTP", "addr", ":3000")
	if err := http.Serve(httpListener, httpHandler); err != nil {
		panic(err)
	}
}
