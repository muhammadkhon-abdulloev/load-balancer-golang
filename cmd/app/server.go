package app

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Proxy    Proxy     `json:"proxy"`
	Services []Service `json"services"`
}

// Proxy is our proxy servers struct (load balancer).
type Proxy struct {
	Port string `json:"port"`
}

// Service is server
type Service struct {
	URL    string `json:"url"`
	IsDead bool
	mx     sync.RWMutex
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

var mx sync.Mutex
var idx int = 0

func (c *Config) Lb(w http.ResponseWriter, r *http.Request) {
	maxLen := len(c.Services)

	mx.Lock()
	currentService := c.Services[idx%maxLen]
	if currentService.GetIsDead() {
		idx++
	}
	targetURL, err := url.Parse(c.Services[idx%maxLen].URL)
	if err != nil {
		log.Fatal(err.Error())
	}
	idx++
	mx.Unlock()
	rp := httputil.NewSingleHostReverseProxy(targetURL)
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("%v is dead", targetURL)
		currentService.SetDead(true)
		c.Lb(w, r)
	}
	rp.ServeHTTP(w, r)
}

func isAlive(url *url.URL) bool {
	conn, err := net.DialTimeout("tcp", url.Host, time.Second*10)
	if err != nil {
		log.Printf("Can't reach to %v, error: ", err.Error())
		return false
	}
	defer conn.Close()
	return true
}

func (c *Config) healthCheck() {
	t := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-t.C:
			for _, service := range c.Services {
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
				log.Printf("service %v checked. status %v", service.URL, msg)
			}
		}
	}
}

func (c *Config) Serve() {

	go c.healthCheck()

	s := http.Server{
		Addr:    ":" + c.Proxy.Port,
		Handler: http.HandlerFunc(c.Lb),
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
