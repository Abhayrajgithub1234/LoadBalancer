package backend

import (
	"sync"
)

type Server struct {
	URL   string
	Alive bool
	mu    sync.RWMutex
}

func (s *Server) SetAlive(status bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Alive = status
}

func (s *Server) IsAlive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Alive
}
