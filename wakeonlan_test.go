package main

import (
	"bytes"
	"net"
	"testing"
)

func TestNew_Only48BitAddresses(t *testing.T) {
	tc := []struct {
		addr    string
		success bool
	}{
		{"00:00:5e:00:53:01", true},
		{"02:00:5e:10:00:00:00:01", false},
		{"00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01", false},
		{"00-00-5e-00-53-01", true},
		{"02-00-5e-10-00-00-00-01", false},
		{"00-00-00-00-fe-80-00-00-00-00-00-00-02-00-5e-10-00-00-00-01", false},
		{"0000.5e00.5301", true},
		{"0200.5e10.0000.0001", false},
		{"0000.0000.fe80.0000.0000.0000.0200.5e10.0000.0001", false},
	}

	for _, tt := range tc {
		t.Run(tt.addr, func(t *testing.T) {
			_, err := New(tt.addr)
			if (err == nil) != tt.success {
				t.Errorf("parsing %+v should fail, but it didnt", tt.addr)
			}
		})
	}
}

func TestNew_Payload(t *testing.T) {
	// get a buffer to hold the packet data
	want := make([]byte, 102)

	// write 6 byte preamble
	copy(want[:6], []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	// write 16 repetitions of target mac address
	offset := 6
	addr, _ := net.ParseMAC("00:00:5e:00:53:01")
	for i := 0; i < 16; i++ {
		copy(want[offset+i*6:offset+i*6+6], addr)
	}

	// construct the magic packet
	var buf bytes.Buffer
	mp, _ := New("00:00:5e:00:53:01")
	mp.Broadcast(&buf)

	// compare!
	got := buf.Bytes()
	if bytes.Compare(want, got) != 0 {
		t.Errorf("want %+v, but got %+v", want, got)
	}
}
