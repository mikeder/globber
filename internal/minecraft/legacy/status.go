package legacy

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Server represents a Minecraft server.
type Server struct {
	Address        string `json:"address"`
	Online         bool   `json:"online"`
	Version        string `json:"version"`
	MOTD           string `json:"motd"`
	CurrentPlayers int    `json:"current_players"`
	MaxPlayers     int    `json:"max_players"`
	Latency        int64  `json:"latency_ms"`
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
		s.Latency = math.MaxInt64
		return errors.Wrap(err, "connecting to server")
	}
	defer conn.Close()
	s.Latency = time.Since(start).Round(time.Millisecond).Nanoseconds()

	return nil
}

// Status updates a server status or returns connection error.
func (s *Server) Status() error {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", s.Address, time.Duration(3)*time.Second)
	if err != nil {
		s.Online = false
		return errors.Wrap(err, "connect to server")
	}
	defer conn.Close()
	s.Latency = time.Since(start).Round(time.Millisecond).Nanoseconds()

	_, err = conn.Write([]byte("\xFE\x01"))
	if err != nil {
		s.Online = false
		return errors.Wrap(err, "write connection")
	}

	data := make([]byte, 512)
	_, err = conn.Read(data)
	if err != nil {
		s.Online = false
		return errors.Wrap(err, "read connection")
	}

	if data == nil || len(data) == 0 {
		s.Online = false
		return nil
	}

	parsed := bytes.Split(data[:], []byte("\x00\x00\x00"))
	if parsed != nil && len(parsed) >= 6 {
		s.Online = true
		s.Version = replaceNulls(parsed[2])
		s.MOTD = replaceNulls(parsed[3])
		s.CurrentPlayers = convertPlayerCount(parsed[4])
		s.MaxPlayers = convertPlayerCount(parsed[5])
	} else {
		s.Online = false
	}

	return nil
}

func replaceNulls(b []byte) string {
	return strings.ReplaceAll(string(b), "\u0000", "")
}

func convertPlayerCount(b []byte) int {
	i, _ := strconv.Atoi(replaceNulls(b))
	return i
}
