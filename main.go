package main

import (
	"context"
	"log"
	"os"
	"time"
)

func main() {
	ctxTO, cancelTO := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelTO()

	ctx, cancel := WithSignal(context.Background(), os.Interrupt)
	defer cancel()

	select {
	case <-ctxTO.Done():
		log.Println("timeout")
	case <-ctx.Done():
		log.Println("signal")
	}
}
