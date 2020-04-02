package minecraft

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/chat"
	_ "github.com/Tnze/go-mc/data/lang/en-us"
	"github.com/google/uuid"
	"github.com/mikeder/globber/internal/models"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// Server represents a Minecraft server.
type Server struct {
	Address        string   `json:"address"`
	Port           int      `json:"port"`
	Online         bool     `json:"online"`
	Version        string   `json:"version"`
	Protocol       int      `json:"protocol"`
	MOTD           string   `json:"motd"`
	CurrentPlayers int      `json:"current_players"`
	MaxPlayers     int      `json:"max_players"`
	OnlinePlayers  []Player `json:"online_players"`
	Latency        int64    `json:"latency_ms"`

	playerDB *sql.DB
}

type Player struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// NewServer returns a server.
func NewServer(addr string, port int, db *sql.DB) *Server {
	srv := &Server{
		Address:  addr,
		Port:     port,
		playerDB: db,
	}

	go srv.periodicUpdate()

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
	s.Online = true
	s.CurrentPlayers = stat.Players.Online
	s.MaxPlayers = stat.Players.Max
	s.Latency = delay.Milliseconds()
	s.OnlinePlayers = stat.Players.Sample
	s.Version = stat.Version.Name
	s.Protocol = stat.Version.Protocol

	go s.updatePlayerTable(s.OnlinePlayers)

	return nil
}

func (s *Server) periodicUpdate() {
	ticker := time.NewTicker(1 * time.Minute)

	for range ticker.C {
		err := s.PingList()
		if err != nil {
			log.Println(errors.Wrap(err, "periodic update"))
		}
	}
}

func (s *Server) updatePlayerTable(players []Player) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*10))
	defer cancel()
	for _, online := range players {
		existing, err := models.PlayerByUUID(ctx, s.playerDB, online.ID.String())
		if err != nil && err != sql.ErrNoRows {
			continue
		}
		if existing == nil {
			newPlayer := models.Player{
				Name:      online.Name,
				UUID:      online.ID.String(),
				FirstSeen: time.Now(),
				LastSeen:  time.Now(),
			}
			err := newPlayer.Save(ctx, s.playerDB)
			if err != nil {
				log.Println(errors.Wrap(err, "add new player"))
			}
		} else {
			existing.LastSeen = time.Now()
			err := existing.Save(ctx, s.playerDB)
			if err != nil {
				log.Println(errors.Wrap(err, "update existing player"))
			}
		}
	}
}
