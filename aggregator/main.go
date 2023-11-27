package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/raphaelmb/go-toll-calculator/types"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	grpcAddr := os.Getenv("AGGREGATOR_GRPC_ENDPOINT")
	httpAddr := os.Getenv("AGGREGATOR_HTTP_ENDPOINT")

	store := makeStore()
	svc := NewInvoiceAggregator(store)
	svc = NewMetricsMiddleware(svc)
	svc = NewLogMiddleware(svc)

	go func() {
		log.Fatal(makeGRPCTransport(grpcAddr, svc))
	}()
	log.Fatal(makeHTTPTransport(httpAddr, svc))
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("GRPC transport running on port", listenAddr)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	server := grpc.NewServer()
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	var (
		aggMetricHandler = NewHTTPMetricHandler("aggregate")
		invMetricHandler = NewHTTPMetricHandler("invoice")
		aggHandler       = makeHTTPHandlerFunc(aggMetricHandler.instrument(handleAggregate(svc)))
		invoiceHandler   = makeHTTPHandlerFunc(invMetricHandler.instrument(handleGetInvoice(svc)))
	)
	http.HandleFunc("/invoice", invoiceHandler)
	http.HandleFunc("/aggregate", aggHandler)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("HTTP transport running on port", listenAddr)
	return http.ListenAndServe(listenAddr, nil)
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
