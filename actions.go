package main

import (
	"log"
	"time"
)

const (
	MEETING_STATUS_CHANGE uint16 = 1
	ACTION_PING           uint16 = 10
)

var actions = map[uint16]func(c *connection, packet Packet){
	ACTION_PING: ActionPing,
}

func ActionPing(c *connection, packet Packet) {
	c.sendBinary <- CreatePacket(ACTION_PING, map[string]string{
		"time": string(time.Now().String()),
	})
}

func Execute(c *connection, packet Packet) {

	if fn, ok := actions[packet.Action]; ok {
		fn(c, packet)
		return
	}

	log.Print("Invalid packet", packet.String())
}
