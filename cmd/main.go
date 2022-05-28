package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	// "os"
	"sync"

	"github.com/muhammadkhon-abdulloev/load-balancer-go/cmd/app"
)

func main() {
	// path := os.Getenv("CONFIG_PATH")

	cfg := new(app.LoadBalancer)
	data, err := ioutil.ReadFile("./configs/config.json")
	if err != nil {
		log.Fatalf(err.Error())
	}
	json.Unmarshal(data, &cfg)

	wg := new(sync.WaitGroup)
	
	wg.Add(1)
	go func ()  {
		cfg.Serve()
		wg.Done()	
	}()


	wg.Wait()
}

