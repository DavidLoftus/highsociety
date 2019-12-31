package http

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"log"
	"sync"
)

type Player struct {
	conn  *websocket.Conn
	mut   sync.RWMutex
	name  string
	lobby *gameLobby
}

func NewPlayer(conn *websocket.Conn) *Player {
	return &Player{
		conn: conn,
	}
}

func (p *Player) SetName(name string) error {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.name = name
	return nil
}

func (p *Player) readPacket() (Packet, error) {
	var anyPacket AnyPacket
	if err := p.conn.ReadJSON(&anyPacket); err != nil {
		log.Println("error reading JSON packet: ", err)
		return nil, err
	}
	return anyPacket.Packet, nil
}

func (p *Player) writePacket(packet Packet) error {
	return p.conn.WriteJSON(AnyPacket{packet})
}

func (p *Player) handlePacket(packet Packet) error {
	switch packet := packet.(type) {
	case NewGamePacket:
		if p.lobby != nil {
			return fmt.Errorf("already in lobby")
		}
		lobby, err := globalGameLobbyStore.NewLobby()
		if err != nil {
			return errors.Wrap(err, "couldn't create new game")
		}
		err = lobby.AddPlayer(p)
		return errors.Wrap(err, "uh oh, failed to add player to his own game")
	case JoinGamePacket:

	case ChangeNamePacket:
		return p.SetName(packet.NewName)
	default:
		return fmt.Errorf("unrecognized packet %T", packet)
	}
	return nil
}

func (p *Player) sendError(err error) error {
	return p.writePacket(&ErrorReportPacket{
		Msg: err.Error(),
	})
}

func (p *Player) Handle() (err error) {
	defer func() {
		closeErr := p.conn.Close()
		if closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				log.Println("Ignored connection close error due to pre-existing error: ", err)
			}
		}
	}()

	for {
		packet, err := p.readPacket()
		if err != nil {
			return errors.Wrap(err, "failed to read packet from websocket")
		}

		err = p.handlePacket(packet)
		if err != nil {
			log.Printf("Error in handling packet from client %s: %v\n", p.conn.LocalAddr(), err)
			if err := p.sendError(err); err != nil {
				return errors.Wrap(err, "failed to report error to client")
			}
		}
	}
}
