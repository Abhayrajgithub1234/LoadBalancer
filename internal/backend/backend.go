package backend


type Server struct {
	URL   string
	Alive bool
}

func (s Server) IsAlive() bool {
	return s.Alive

}

func (s *Server) SetAlive(status bool) {
	s.Alive = status 
}

