package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	usage = `
Usage: wakeonlan [options] macaddr

  Sends a specially designed Ethernet frame that "awakens" a networked computer.

Options:

  -b=<address>
    The broadcast address to use. Defaults to 255.255.255.255

  -p=<port>
    The port to use for the UDP datagram. Typically sent to port 0, 7 or 9.
    Defaults to port 9 which maps to the well-known discard protocol.
`
)

func main() {
	os.Exit(Run(os.Stdin, os.Stdout, os.Stdout, os.Args[1:]))
}

func Run(stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	var broadcast, port string

	flags := flag.NewFlagSet("wakeonlan", flag.ContinueOnError)
	flags.StringVar(&broadcast, "b", "", "address")
	flags.StringVar(&port, "p", "", "port")
	flags.Usage = func() {
		fmt.Fprintln(stderr, strings.TrimSpace(usage))
	}

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if broadcast == "" {
		broadcast = "255.255.255.255"
	}

	if port == "" {
		port = "9"
	}

	arguments := flags.Args()
	if len(arguments) != 1 {
		fmt.Fprintln(stderr, strings.TrimSpace(usage))
		return 1
	}

	address := fmt.Sprintf("%s:%s", broadcast, port)
	conn, err := net.Dial("udp", address)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	defer conn.Close()

	macaddr := strings.TrimSpace(arguments[0])
	mp, err := New(macaddr)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	err = mp.Broadcast(conn)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	fmt.Fprintln(stdout, "Magic packet sent to", macaddr)
	return 0
}
