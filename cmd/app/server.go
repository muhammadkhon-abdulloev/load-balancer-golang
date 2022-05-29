package app

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/muhammadkhon-abdulloev/load-balancer-go/pkg/service"
)

type LoadBalancer struct {
	httpServer *http.Server
	Port       string `json:"port"`
	Services   []service.Service
	currentSvc int
	sync.RWMutex
}

func (lb *LoadBalancer) Lb(w http.ResponseWriter, r *http.Request) {
	maxLen := len(lb.Services)

	lb.RLock()
	currentService := lb.Services[lb.currentSvc%maxLen]
	if currentService.IsDead() {
		lb.currentSvc++
	}
	targetURL, err := url.Parse(lb.Services[lb.currentSvc%maxLen].URL)
	if err != nil {
		log.Fatal(err.Error())
	}
	lb.currentSvc++
	lb.RUnlock()
	rp := httputil.NewSingleHostReverseProxy(targetURL)
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("%v is dead", targetURL)
		currentService.SetDead(true)
		lb.Lb(w, r)
	}
	rp.ServeHTTP(w, r)
}

func (lb *LoadBalancer) healthCheck() {
	t := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-t.C:
			for _, svc := range lb.Services {
				pingURL, err := url.Parse(svc.URL)
				if err != nil {
					log.Fatal(err.Error())
				}
				isAlive := service.IsAlive(pingURL)
				svc.SetDead(!isAlive)
				msg := "ok"
				if !isAlive {
					msg = "dead"
				}
				log.Printf("service %v checked. status %v", svc.URL, msg)
			}
		}
	}
}

func (lb *LoadBalancer) Serve() error {

	go lb.healthCheck()

	lb.httpServer = &http.Server{
		Addr:    ":" + lb.Port,
		Handler: http.HandlerFunc(lb.Lb),
	}
	return lb.httpServer.ListenAndServe()

}

func (lb *LoadBalancer) Shutdown(ctx context.Context) error {
	return lb.httpServer.Shutdown(ctx)
}
