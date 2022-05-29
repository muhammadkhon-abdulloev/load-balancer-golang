package service

import (
	"log"
	"net"
	"net/url"
	"sync"
	"time"
)

// Service is target service structure. URL is url of target service. isDead, is status of service
//(If isDead = false, service is alive. Neither, service is die)
type Service struct {
	URL    string `json:"url"`
	isDead bool
	mx     sync.RWMutex
}

// SetDead - is the method of the Service structure which sets the service status. Argument - boolean type. return nothing
func (s *Service) SetDead(status bool) {
	s.mx.Lock()
	s.isDead = status
	s.mx.Unlock()
}

// IsDead - is the method of the Service structure, which returns status of service. 
// Returns false when service is alive and true otherwise
func (s *Service) IsDead() bool {
	s.mx.RLock()
	status := s.isDead
	s.mx.RUnlock()
	return status
}


// IsAlive - is the function which returns status of gived url.
// Return true when service is alive and false otherwise
func IsAlive(url *url.URL) bool {
	conn, err := net.DialTimeout("tcp", url.Host, time.Millisecond * 500)
	if err != nil {
		log.Printf("Can't reach to %v, error: ", err.Error())
		return false
	}
	defer conn.Close()
	return true
}