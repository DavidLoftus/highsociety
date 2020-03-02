package http

import (
	"encoding/json"
	"fmt"
)

type PacketType string

const (
	// Client Packets
	PacketNewGame    PacketType = "NEW_GAME"
	PacketJoinGame              = "JOIN_GAME"
	PacketChangeName            = "CHANGE_NAME"

	// Server Packets
	PacketOk           = "OK"
	PacketErrorReport  = "ERROR_REPORT"
	PacketPlayerStatus = "PLAYER_STATUS"

	PacketArbitrary = "ARBITRARY"
)

type PacketHeader struct {
	Id int `json:"id"`
}

func (p *PacketHeader) Header() *PacketHeader {
	return p
}

type Packet interface {
	GetType() PacketType
	Header() *PacketHeader
}

type AnyPacket struct {
	Packet
}

type packetWrapper struct {
	Type PacketType  `json:"type"`
	Msg  interface{} `json:"msg"`
}

func (a *AnyPacket) UnmarshalJSON(data []byte) error {
	// RawMessage to store packet contents
	var msg json.RawMessage

	// Wrapper to get type aswell as raw message
	wrapper := packetWrapper{
		Msg: &msg,
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return err
	}

	// Create instance of whichever packet type it is
	var packet Packet
	switch wrapper.Type {
	case PacketNewGame:
		packet = new(NewGamePacket)
	case PacketJoinGame:
		packet = new(JoinGamePacket)
	case PacketChangeName:
		packet = new(ChangeNamePacket)
	case PacketPlayerStatus:
		packet = new(PlayerStatusPacket)
	case PacketOk:
		packet = new(OkPacket)
	default:
		return fmt.Errorf("unrecognised packet type: %q", wrapper.Type)
	}

	// Attempt to unmarshal RawMessage to selected packet type
	if err := json.Unmarshal(msg, packet); err != nil {
		return err
	}
	a.Packet = packet

	return nil
}

func (a AnyPacket) MarshalJSON() ([]byte, error) {
	wrapper := packetWrapper{
		Type: a.GetType(),
		Msg:  a.Packet,
	}

	return json.Marshal(wrapper)
}

type NewGamePacket struct {
	PacketHeader
}

func (NewGamePacket) GetType() PacketType {
	return PacketNewGame
}

type JoinGamePacket struct {
	PacketHeader
	LobbyID string `json:"lobby_id"`
}

func (JoinGamePacket) GetType() PacketType {
	return PacketJoinGame
}

type ChangeNamePacket struct {
	PacketHeader
	NewName string `json:"new_name"`
}

func (ChangeNamePacket) GetType() PacketType {
	return PacketChangeName
}

type ErrorReportPacket struct {
	PacketHeader
	Msg string `json:"msg"`
}

func (ErrorReportPacket) GetType() PacketType {
	return PacketErrorReport
}

type PlayerStatusPacket struct {
	PacketHeader
	Name string `json:"name"`
}

func (PlayerStatusPacket) GetType() PacketType {
	return PacketPlayerStatus
}

type OkPacket struct {
	PacketHeader
}

func (OkPacket) GetType() PacketType {
	return PacketOk
}

type ArbitraryPacket struct {
	PacketHeader
	Data map[string]interface{} `json:"data"`
}

func (ArbitraryPacket) GetType() PacketType {
	return PacketArbitrary
}
