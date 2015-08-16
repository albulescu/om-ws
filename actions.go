package main

import (
	"log"
	"time"
)

const (
	MEETING_STATUS_CHANGE uint16 = 1
	ACTION_PING           uint16 = 10
)

type Action func(c *connection, data PacketData)

var actions = map[uint16]Action{
	ACTION_PING: ActionPing,
}

func Auth(c *connection, data PacketData) {

}

func ActionPing(c *connection, data PacketData) {
	c.sendBinary <- CreatePacket(ACTION_PING, map[string]string{
		"time": string(time.Now().String()),
	})
}

func Execute(c *connection, packet Packet) {

	if fn, ok := actions[packet.Action]; ok {
		fn(c, packet.Data)
		return
	}

	log.Print("Invalid packet", packet.String())
}
