package main

import (
	"context"
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/muhammadkhon-abdulloev/load-balancer-go/cmd/app"
)

func main() {
	// path := os.Getenv("CONFIG_PATH")

	
	data, err := ioutil.ReadFile("./configs/config.json")
	if err != nil {
		log.Fatalf(err.Error())
	}

	lb := new(app.LoadBalancer)
	json.Unmarshal(data, &lb)

	wg := new(sync.WaitGroup)
	
	wg.Add(1)
	go func ()  {
		
		if err := lb.Serve(); err != nil {
			log.Fatalf("error occured while running server")
		}
		
		wg.Done()	
	}()

		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
		<- exit
		
		log.Println("Shutting down proxy server...")

		if err := lb.Shutdown(context.Background()); err != nil {
			log.Fatalf("error occured on server while shutting down: %s", err.Error())
		}

	wg.Wait()
}

