package rcon

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/tankbusta/renx-rcon/events"
	"github.com/tankbusta/renx-rcon/games"
)

const WriterSizeQueue = 10

type Server struct {
	// Game indicates which Totem Arts game this server is running
	Game games.Game

	// Address of the Game Server
	Address string

	// ConnectionID of the RCON connection
	ConnectionID string

	// Version of the UDK Game Server
	Version events.Version

	// IsAuthenticated marks if the server has accepted our RCON credentials
	IsAuthenticated bool

	// IsConnected marks if the server has been connected to (but not necessairly authenticated)
	IsConnected bool

	GameState *GameStateManager

	// unexported fields below
	rconPassword  string
	pendingWrites chan string
}

func NewServer(rconPassword, gameServer string) *Server {
	writes := make(chan string, WriterSizeQueue)
	return &Server{
		Address:       gameServer,
		rconPassword:  rconPassword,
		pendingWrites: writes,
	}
}

// Connect to the UDK game server and authenticate with the RCON password
func (s *Server) Connect(ctx context.Context) (net.Conn, error) {
	var d net.Dialer

	s.IsConnected = false
	s.IsAuthenticated = false

	conn, err := d.DialContext(ctx, "tcp", s.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RCON at %s: %w", s.Address, err)
	}

	if _, err := conn.Write([]byte(fmt.Sprintf("a%s\n", s.rconPassword))); err != nil {
		return nil, fmt.Errorf("failed to authenticate to RCON at %s: %w", s.Address, err)
	}

	s.IsConnected = true
	return conn, nil
}

// Ready indicates this server has been connected to and authentication acknowledged
func (s *Server) Ready() bool { return s.IsConnected && s.IsAuthenticated }

// Destroy should be called when we no longer need this server.
// It frees up resources. Once destroyed, the server should be re-created to avoid issues
func (s *Server) Destroy() {
	s.IsConnected = false
	s.IsAuthenticated = false

	close(s.pendingWrites)
	s.pendingWrites = nil
}

func (s *Server) WriteMsg(msg string) {
	s.pendingWrites <- msg
}

func (s *Server) Start(ctx context.Context) error {
	state := NewGameState(s)

	go func() {
		state.Start(ctx)
	}()

MainLoop:
	for {
		select {
		case <-ctx.Done():
			break MainLoop
		default:
			var rdr bufio.Reader

			{
				conn, err := s.Connect(ctx)
				if err != nil {
					return err
				}
				defer conn.Close()

				rdr.Reset(conn)

			ReadLoop:
				for {
					select {
					case <-ctx.Done():
						break MainLoop
					// Design Note: Writes to RCON are so in-frequent here
					// we're going to use the same loop for both reading and writing
					case writeMe := <-s.pendingWrites:
						log.Printf("[ !! ] Writing message to rcon: %s\n", writeMe)
						conn.SetWriteDeadline(time.Now().Add(time.Second * 2))
						if _, err := conn.Write([]byte(writeMe)); err != nil {
							log.Printf("Failed to write to RCON msg at %s: %s", s.Address, err)
						}
						// TODO: We could put it back in the channel if it fails?
					default:
					}

					conn.SetReadDeadline(time.Now().Add(time.Second * 1))
					msg, err := rdr.ReadString('\n')
					if err != nil {
						if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
							continue ReadLoop // Keep going
						}

						log.Printf("Failed to read RCON msg: %s\n", err)
						break ReadLoop
					}

					if len(msg) < 2 {
						log.Printf("RCON msg length of %d too small\n", len(msg))
						break ReadLoop
					}

					_ = msg[1] // Bounds check

					msgNoType := msg[1:]
					switch events.ServerType(msg[0]) {
					case events.RCONGameVersion:
						var ver events.Version

						if err := ver.Parse(msgNoType); err != nil {
							// If we cant parse the version, we're gonna bomb out
							// because we might run into unexpected behavior
							return err
						}

						switch {
						case (ver.GameVersion > 12000 && ver.GameVersion < 13000):
							s.Game = games.GameRenegadeX
						default:
							s.Game = games.GameUnknown
						}

						log.Printf("[ !! ] %s", ver)
					case events.AuthenticationSuccess:
						s.ConnectionID = msgNoType
						s.IsAuthenticated = true

						log.Printf("[ ++ ] Got AuthSuccess, starting event stream!\n")
						// Authenticated and ready to accept streaming!
						if _, err := conn.Write([]byte(fmt.Sprintf("s\n"))); err != nil {
							log.Fatalf("Failed to write RCON auth: %s\n", err)
						}
					case events.Error:
						var err events.ServerError
						err.Parse(msgNoType)

						// If we're not authenticated and we get an error, bomb out
						if !s.IsAuthenticated {
							return err
						}

						// Otherwise, just log the error
						log.Printf("[ XX ] RCON error: %s\n", err)
					case events.CommandResponse:
						fmt.Println(hex.Dump([]byte(msgNoType)))
						fallthrough
					case events.CommandExecutionFinished:
						fallthrough
					case events.GameLog:
						fallthrough
					case events.ServerDevBot:
						fmt.Println(msg)
					}
				}
			}
		}
	}

	s.IsConnected = false
	log.Printf("[ !! ] Goodbye!")
	return nil
}
