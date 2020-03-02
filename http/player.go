package http

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"sync"
)

type Player struct {
	conn  *websocket.Conn
	mut   sync.RWMutex
	name  string
	lobby *gameLobby
}

func NewPlayer(conn *websocket.Conn) *Player {
	player := &Player{
		conn: conn,
		name: fmt.Sprintf("Guest#%04d", rand.Intn(10000)),
	}
	//player.sendStatus()
	return player
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

func (p *Player) handlePacket(packet Packet) (Packet, error) {
	switch packet := packet.(type) {
	case *NewGamePacket:
		if p.lobby != nil {
			return nil, fmt.Errorf("already in lobby")
		}
		lobby, err := globalGameLobbyStore.NewLobby()
		if err != nil {
			return nil, errors.Wrap(err, "couldn't create new game")
		}
		if err := lobby.AddPlayer(p); err != nil {
			return nil, errors.Wrap(err, "uh oh, failed to add player to his own game")
		}
		resp := &ArbitraryPacket{
			Data: map[string]interface{}{
				"lobby_id": lobby.id,
			},
		}
		return resp, nil
	case *JoinGamePacket:
		if p.lobby != nil {
			return nil, fmt.Errorf("already in lobby")
		}
		lobby := globalGameLobbyStore.GetLobby(packet.LobbyID)
		if lobby == nil {
			return nil, fmt.Errorf("no such lobby exists")
		}

		if err := lobby.AddPlayer(p); err != nil {
			return nil, err
		}

		resp := &ArbitraryPacket{
			Data: map[string]interface{}{
				"lobby_id": lobby.id,
			},
		}
		return resp, nil
	case *ChangeNamePacket:
		return nil, p.SetName(packet.NewName)
	default:
		return nil, fmt.Errorf("unrecognized packet %T", packet)
	}
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

		reqHeader := packet.Header()

		response, err := p.handlePacket(packet)
		if err != nil {
			log.Printf("Error in handling packet from client %s: %v\n", p.conn.LocalAddr(), err)

			// We should send an error packet
			response = &ErrorReportPacket{Msg: err.Error()}
		}

		if response == nil {
			response = &OkPacket{}
		}

		respHeader := response.Header()
		respHeader.Id = reqHeader.Id

		if err := p.writePacket(response); err != nil {
			return errors.Wrap(err, "failed to send packet to client")
		}
	}
}
