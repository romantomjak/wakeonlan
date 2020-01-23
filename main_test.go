package main

import (
	"bytes"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"
)

func assertEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func assertContains(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("expected %q to contain %q, but it didn't", got, want)
	}
}

func TestRun_RequiresMacAddress(t *testing.T) {
	stdin := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	args := []string{}

	code := Run(stdin, stdout, stderr, args)
	assertEqual(t, code, 1)

	out := stderr.String()
	assertContains(t, out, "Usage:")
}

func TestRun_CanSendWakeOnLan(t *testing.T) {
	stdin := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	args := []string{"-b=127.0.0.1", "-p=3233", "00:00:5e:00:53:01"}

	go func() {
		Run(stdin, stdout, stderr, args)
	}()

	conn, err := net.ListenPacket("udp", "127.0.0.1:3233")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(time.Second)); err != nil {
		t.Error(err)
	}

	buf := make([]byte, 102)
	if _, _, err := conn.ReadFrom(buf); err != nil {
		t.Error(err)
	}

	out := stdout.String()
	assertContains(t, out, "Magic packet sent to 00:00:5e:00:53:01")
}
