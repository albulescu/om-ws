package main

import (
	"bytes"
	"encoding/binary"
	"time"
)

var Packet struct {
	Event string
	Time  time.Time
	Data  interface{}
}

var actions = map[string]func(c *connection){
	"ping": ActionPing,
}

func CreatePacket(action int, body map[string]string) []byte {

	buffer := new(bytes.Buffer)
	separator := []byte{0xc0, 0x80}

	binary.Write(buffer, binary.BigEndian, []byte{'o', 'm'}) //control
	binary.Write(buffer, binary.BigEndian, uint8(1))         //version
	binary.Write(buffer, binary.BigEndian, uint8(action))    //action

	var size int = 0
	for key, value := range body {
		size += (len(key) + 2 + len(value) + 2)
	}

	binary.Write(buffer, binary.BigEndian, uint32(size)) //body size

	for key, value := range body {
		binary.Write(buffer, binary.BigEndian, []byte(key))
		binary.Write(buffer, binary.BigEndian, separator)
		binary.Write(buffer, binary.BigEndian, []byte(value))
		binary.Write(buffer, binary.BigEndian, separator)
	}

	return buffer.Bytes()
}

func ActionPing(c *connection) {
	c.sendBinary <- CreatePacket(10, map[string]string{
		"time": string(time.Now().String()),
	})
}

func Execute(event string, data interface{}) {

}
