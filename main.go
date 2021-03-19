package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	stopChan := make(chan struct{})

	ls, err := net.Listen("tcp", ":7319")
	if err != nil {
		log.Fatalf("%v", err)
	}

	srv := &http.Server{
		Handler:      new(ScaleHandler),
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(stopChan)
	}()

	err = srv.Serve(ls)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("%v", err)
	}

	<-stopChan
}

type ScaleHandler struct {
}

func (sh *ScaleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v", ScaleResult{Error: false, Weight: 1.1234})
}

type ScaleResult struct {
	Error            bool
	ErrorDescription string
	Weight           float32
}

func (sr ScaleResult) String() string {
	if sr.Error {
		return fmt.Sprintf(`{"error":true, "error_description":%s}`, sr.ErrorDescription)
	}

	return fmt.Sprintf(`{"error":false, "error_description":%f}`, sr.Weight)
}
