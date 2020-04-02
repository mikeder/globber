package minecraft

import (
	"encoding/json"
	"time"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/chat"
	_ "github.com/Tnze/go-mc/data/lang/en-us"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Server represents a Minecraft server.
type Server struct {
	Address        string        `json:"address"`
	Port           int           `json:"port"`
	Online         bool          `json:"online"`
	Version        string        `json:"version"`
	Protocol       int           `json:"protocol"`
	MOTD           string        `json:"motd"`
	CurrentPlayers int           `json:"current_players"`
	MaxPlayers     int           `json:"max_players"`
	OnlinePlayers  []Player      `json:"online_players"`
	Latency        time.Duration `json:"latency"`
}

type Player struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// NewServer returns a server.
func NewServer(addr string, port int) *Server {
	srv := &Server{
		Address: addr,
		Port:    port,
	}
	return srv
}

type status struct {
	Description chat.Message
	Players     struct {
		Max    int
		Online int
		Sample []Player
	}
	Version struct {
		Name     string
		Protocol int
	}
	//favicon ignored
}

// PingList performs a ping and list player command on the server.
func (s *Server) PingList() error {
	resp, delay, err := bot.PingAndList(s.Address, s.Port)
	if err != nil {
		s.Online = false
		return errors.Wrap(err, "ping and list players")
	}

	var stat status
	err = json.Unmarshal(resp, &stat)
	if err != nil {
		return errors.Wrap(err, "unmarshal status")
	}

	s.MOTD = stat.Description.String()
	s.CurrentPlayers = stat.Players.Online
	s.MaxPlayers = stat.Players.Max
	s.Latency = delay
	s.OnlinePlayers = stat.Players.Sample
	s.Version = stat.Version.Name
	s.Protocol = stat.Version.Protocol

	return nil
}
