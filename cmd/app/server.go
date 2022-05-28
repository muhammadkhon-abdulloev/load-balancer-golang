package app

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
	"log"
)

type LoadBalancer struct {
	Port string `json:"port"`
	Services []Service
	current int
	sync.Mutex
}


// Service is server
type Service struct {
	URL    string `json:"url"`
	IsDead bool
	mx     sync.RWMutex
}

func NewLB(port string) *LoadBalancer{
	return &LoadBalancer{
		Port: port,
	}
}

func (s *Service) SetDead(status bool) {
	s.mx.Lock()
	s.IsDead = status
	s.mx.Unlock()
}

func (s *Service) GetIsDead() bool {
	s.mx.RLock()
	status := s.IsDead
	s.mx.RUnlock()
	return status
}

func (lb *LoadBalancer) Lb(w http.ResponseWriter, r *http.Request) {
	maxLen := len(lb.Services)

	lb.Lock()
	currentService := lb.Services[lb.current%maxLen]
	if currentService.GetIsDead() {
		lb.current++
	}
	targetURL, err := url.Parse(lb.Services[lb.current%maxLen].URL)
	if err != nil {
		log.Fatal(err.Error())
	}
	lb.current++
	lb.Unlock()
	rp := httputil.NewSingleHostReverseProxy(targetURL)
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("%v is dead", targetURL)
		currentService.SetDead(true)
		lb.Lb(w, r)
	}
	rp.ServeHTTP(w, r)
}

func isAlive(url *url.URL) bool {
	conn, err := net.DialTimeout("tcp", url.Host, time.Millisecond * 500)
	if err != nil {
		log.Printf("Can't reach to %v, error: ", err.Error())
		return false
	}
	defer conn.Close()
	return true
}

func (lb *LoadBalancer) healthCheck() {
	t := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-t.C:
			for _, service := range lb.Services {
				pingURL, err := url.Parse(service.URL)
				if err != nil {
					log.Fatal(err.Error())
				}
				isAlive := isAlive(pingURL)
				service.SetDead(!isAlive)
				msg := "ok"
				if !isAlive {
					msg = "dead"
				}
				lg := fmt.Sprintf("service %v checked. status %v", service.URL, msg)
				log.Printf(lg)
			}
		}
	}
}

func (lb *LoadBalancer) Serve() {

	go lb.healthCheck()

	s := http.Server{
		Addr:    ":" + lb.Port,
		Handler: http.HandlerFunc(lb.Lb),
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}

