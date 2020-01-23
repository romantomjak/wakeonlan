package main

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

// A HardwareAddr represents a physical hardware address
type HardwareAddr [6]byte

// MagicPacket is a broadcast Ethernet frame containing anywhere within its
// payload 6 bytes of all 255 (FF in hex), followed by sixteen repetitions
// of the target computer's 48-bit MAC address
type MagicPacket struct {
	SyncStream [6]byte
	TargeMAC   [16]HardwareAddr
}

// New returns a new MagicPacket
func New(mac string) (*MagicPacket, error) {
	addr, err := net.ParseMAC(mac)
	if err != nil {
		return nil, err
	}

	// accept only 48-bit MAC addresses in various formats
	if len(addr) != 6 {
		return nil, errors.New("invalid 48-bit MAC address")
	}

	p := &MagicPacket{}
	p.SyncStream = [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	for i := 0; i < 16; i++ {
		copy(p.TargeMAC[i][:6], addr)
	}

	return p, nil
}

// Broadcast writes the magic packet to the underlying data stream
func (mp *MagicPacket) Broadcast(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, mp)
}
