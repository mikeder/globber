package minecraft

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Server represents a Minecraft server.
type Server struct {
	Address        string
	Online         bool
	Version        string
	MOTD           string
	CurrentPlayers string
	MaxPlayers     string
	Latency        time.Duration
}

// NewServer returns a server.
func NewServer(host string, port int) *Server {
	addr := fmt.Sprintf("%s:%d", host, port)
	srv := &Server{
		Address: addr,
	}
	return srv
}

// Ping returns a server latency and/or connection error.
func (s *Server) Ping() (time.Duration, error) {
	start := time.Now()
	_, err := net.DialTimeout("tcp", s.Address, time.Duration(3)*time.Second)
	s.Latency = time.Since(start).Round(time.Millisecond)
	return s.Latency, err
}

// Status updates a server status or returns connection error.
func (s *Server) Status() error {
	conn, err := net.DialTimeout("tcp", s.Address, time.Duration(3)*time.Second)
	defer conn.Close()

	_, err = conn.Write([]byte("\xFE\x01"))
	if err != nil {
		s.Online = false
		return err
	}

	data := make([]byte, 512)
	_, err = conn.Read(data)
	if err != nil {
		s.Online = false
		return err
	}

	if data == nil || len(data) == 0 {
		s.Online = false
		return err
	}

	parsed := strings.Split(string(data[:]), "\x00\x00\x00")
	if parsed != nil && len(parsed) >= 6 {
		s.Online = true
		s.Version = parsed[2]
		s.MOTD = parsed[3]
		s.CurrentPlayers = parsed[4]
		s.MaxPlayers = parsed[5]
	} else {
		s.Online = false
	}

	return nil
}
