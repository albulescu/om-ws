package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

var (
	SEPARATOR = []byte{0xc0, 0x80}
)

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

type PacketData struct {
	Map map[string]string
}

func (p *PacketData) Get(key string) (string, error) {

	if p.Has(key) {
		return p.Map[key], nil
	}

	return "", &errorString{fmt.Sprintf("No data with key %s", key)}
}

func (p *PacketData) Set(key string, value string) error {

	if p.Has(key) {
		return &errorString{fmt.Sprintf("Key %s already exist", key)}
	}

	p.Map[key] = value

	return nil
}

func (p *PacketData) Has(key string) bool {
	if _, ok := p.Map[key]; ok {
		return true
	}
	return false
}

func (p *PacketData) GetInt(key string) (int64, error) {

	if p.Has(key) {
		value, _ := p.Get(key)
		return strconv.ParseInt(value, 10, 8)
	}

	return 0, &errorString{fmt.Sprintf("No data with key %s", key)}
}

func (p *PacketData) GetBool(key string) (bool, error) {

	if p.Has(key) {
		value, _ := p.Get(key)
		return (value == "true" || value == "1"), nil
	}

	return false, &errorString{fmt.Sprintf("No data with key %s", key)}
}

type Packet struct {
	Control [2]byte
	Version uint8
	Action  uint16
	Size    uint32
	Data    PacketData
}

func (p *Packet) String() string {
	return fmt.Sprintf("[ Control:%s Version:%s Action:%s BodySize:%s]", p.Control, p.Version, p.Action, p.Size)
}

func (p *Packet) Get(key string) (string, error) {
	return p.Data.Get(key)
}

func ScanBody(data []byte, atEOF bool) (advance int, token []byte, err error) {

	length := len(data)

	for i := 0; i < length; i++ {
		if data[i] == 0xc0 && data[i+1] == 0x80 {
			return i + 2, data[0:i], nil
		}
	}

	return 0, nil, nil
}

func ReadPacket(b []byte) (Packet, error) {

	p := Packet{}
	buf := bytes.NewReader(b)

	binary.Read(buf, binary.BigEndian, &p.Control)
	binary.Read(buf, binary.BigEndian, &p.Version)
	binary.Read(buf, binary.BigEndian, &p.Action)
	binary.Read(buf, binary.BigEndian, &p.Size)

	var body = make([]byte, p.Size)
	binary.Read(buf, binary.BigEndian, &body)

	buf = bytes.NewReader(body)
	reader := bufio.NewReader(buf)
	scanner := bufio.NewScanner(reader)
	scanner.Split(ScanBody)

	toggle := true
	key := ""
	data := make(map[string]string)

	for scanner.Scan() {

		if toggle {
			key = scanner.Text()
		} else {
			data[key] = scanner.Text()
		}

		toggle = !toggle
	}

	p.Data.Map = data

	return p, nil
}

func CreatePacket(action uint16, body map[string]string) []byte {

	buffer := new(bytes.Buffer)
	separator := []byte{0xc0, 0x80}

	binary.Write(buffer, binary.BigEndian, []byte{'o', 'm'}) //control
	binary.Write(buffer, binary.BigEndian, uint8(1))         //version
	binary.Write(buffer, binary.BigEndian, action)           //action

	var size int = 0
	for key, value := range body {
		size += (len(key) + len(value) + 4)
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
