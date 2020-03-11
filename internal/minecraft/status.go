package minecraft

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Server represents a Minecraft server.
type Server struct {
	Address        string        `json:"address"`
	Online         bool          `json:"online"`
	Version        string        `json:"version"`
	MOTD           string        `json:"motd"`
	CurrentPlayers string        `json:"current_players"`
	MaxPlayers     string        `json:"max_players"`
	Latency        time.Duration `json:"latency"`
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
func (s *Server) Ping() error {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", s.Address, time.Duration(3)*time.Second)
	if err != nil {
		return errors.Wrap(err, "connecting to server")
	}
	defer conn.Close()

	s.Latency = time.Since(start).Round(time.Millisecond)
	return nil
}

// Status updates a server status or returns connection error.
func (s *Server) Status() error {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", s.Address, time.Duration(3)*time.Second)
	if err != nil {
		return errors.Wrap(err, "connecting to server")
	}
	defer conn.Close()
	s.Latency = time.Since(start).Round(time.Millisecond)

	_, err = conn.Write([]byte("\xFE\x01"))
	if err != nil {
		s.Online = false
		return errors.Wrap(err, "writing to server connection")
	}

	data := make([]byte, 512)
	_, err = conn.Read(data)
	if err != nil {
		s.Online = false
		return errors.Wrap(err, "reading from server connection")
	}

	if data == nil || len(data) == 0 {
		s.Online = false
		return err
	}

	parsed := bytes.Split(data[:], []byte("\x00\x00\x00"))
	if parsed != nil && len(parsed) >= 6 {
		s.Online = true
		s.Version = replaceNulls(parsed[2])
		s.MOTD = replaceNulls(parsed[3])
		s.CurrentPlayers = replaceNulls(parsed[4])
		s.MaxPlayers = replaceNulls(parsed[5])
	} else {
		s.Online = false
	}

	return nil
}

func replaceNulls(b []byte) string {
	return strings.ReplaceAll(string(b), "\u0000", "")
}
