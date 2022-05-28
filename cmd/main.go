package main

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/muhammadkhon-abdulloev/load-balancer-go/cmd/app"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := new(app.Config)
	data, err := ioutil.ReadFile("./configs/config.json")
	if err != nil {
		log.Fatalf(err.Error())
	}
	json.Unmarshal(data, &cfg)

	wg := new(sync.WaitGroup)
	wg.Add(5)

	go func ()  {
		cfg.Serve()
		wg.Done()	
	}()


	wg.Wait()
}

