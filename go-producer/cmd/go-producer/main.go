package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	var (
		envHTTPAddr = envString("HTTP_ADDR", ":8080")
		httpAddr    = *flag.String("http.addr", envHTTPAddr, "HTTP listen address")
	)
	flag.Parse()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	mux := http.NewServeMux()
	mux.Handle("/_status/liveness", livenessHandler())
	mux.Handle("/_status/readiness", readinessHandler())

	server := &http.Server{Addr: httpAddr, Handler: mux}

	go func(server *http.Server) {
		log.Println("service started.")
		log.Printf("listening to %s\n", httpAddr)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}(server)

	<-signals

	fmt.Println("gracefully shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("service stopped.")

}

func livenessHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
}

func readinessHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
}

func envString(env, fallback string) string {
	e, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	return e
}
