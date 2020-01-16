package http

import (
	"encoding/json"
	"fmt"
)

type PacketType string

const (
	// Client Packets
	PacketNewGame     PacketType = "NEW_GAME"
	PacketJoinGame               = "JOIN_GAME"
	PacketChangeName             = "CHANGE_NAME"
	PacketErrorReport            = "ERROR_REPORT"

	// Server Packets
	PacketPlayerStatus = "PLAYER_STATUS"
)

type Packet interface {
	GetType() PacketType
}

type AnyPacket struct {
	Packet
}

type packetWrapper struct {
	Type PacketType  `json:"type"`
	Msg  interface{} `json:"msg"`
}

func (a AnyPacket) UnmarshalJSON(data []byte) error {
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
}

func (NewGamePacket) GetType() PacketType {
	return PacketNewGame
}

type JoinGamePacket struct {
	LobbyID string
}

func (JoinGamePacket) GetType() PacketType {
	return PacketJoinGame
}

type ChangeNamePacket struct {
	NewName string
}

func (ChangeNamePacket) GetType() PacketType {
	return PacketChangeName
}

type ErrorReportPacket struct {
	Msg string
}

func (ErrorReportPacket) GetType() PacketType {
	return PacketErrorReport
}

type PlayerStatusPacket struct {
	name string
}

func (PlayerStatusPacket) GetType() PacketType {
	return PacketPlayerStatus
}
