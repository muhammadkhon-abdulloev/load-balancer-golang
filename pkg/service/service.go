package service

import (
	"log"
	"net"
	"net/url"
	"sync"
	"time"
)

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

func IsAlive(url *url.URL) bool {
	conn, err := net.DialTimeout("tcp", url.Host, time.Millisecond * 500)
	if err != nil {
		log.Printf("Can't reach to %v, error: ", err.Error())
		return false
	}
	defer conn.Close()
	return true
}