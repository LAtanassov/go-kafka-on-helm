package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

func main() {

	var (
		envHTTPAddr = envString("HTTP_ADDR", ":8080")
		httpAddr    = *flag.String("http.addr", envHTTPAddr, "HTTP listen address")

		envKafkaAddrs = envString("KAFKA_ADDRS", ":9092")
		envKafkaTopic = envString("KAFKA_TOPIC", "timer")

		kafkaAddrs = *flag.String("kafka.addrs", envKafkaAddrs, "comma-seperated kafka broker addresses")
		kafkaTopic = *flag.String("kafka.topic", envKafkaTopic, "kafka topic")
	)
	flag.Parse()

	// os signal trap
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// health status endpoint
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

	// kafka producer
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(strings.Split(kafkaAddrs, ","), config)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			log.Println("message acknowledged.")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			log.Println(err)
		}
	}()

ProducerLoop:
	for {
		select {
		case <-time.After(5 * time.Second):
			producer.Input() <- &sarama.ProducerMessage{Topic: kafkaTopic, Value: sarama.StringEncoder(time.Now().String())}
			log.Println("message sent.")
		case <-signals:
			producer.AsyncClose()
			break ProducerLoop
		}
	}

	fmt.Println("gracefully shutdown...")
	wg.Wait()

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
