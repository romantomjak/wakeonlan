# wake-on-lan

Wake-on-LAN ("WOL") is an Ethernet computer networking standard that allows a computer to be turned on or awakened by a network message. It is implemented using a specially designed Ethernet frame called a magic packet, which is sent to all computers in a network, among them the computer to be awakened.

The magic packet is a broadcast frame containing anywhere within its payload 6 bytes of all 255 (FF FF FF FF FF FF in hexadecimal), followed by sixteen repetitions of the target computer's 48-bit MAC address, for a total of 102 bytes.

## Installation

```sh
go get -u github.com/romantomjak/wakeonlan
```

## Usage

```sh
Usage: wakeonlan [options] macaddr

  Sends a specially designed Ethernet frame that "awakens" a networked computer.

Options:

  -b=<address>
    The broadcast address to use. Defaults to 255.255.255.255

  -p=<port>
    The port to use for the UDP datagram. Typically sent to port 0, 7 or 9.
    Defaults to port 9 which maps to the well-known discard protocol.
```

To awaken a computer:

```sh
$ wakeonlan 00:00:5e:00:53:01
Magic packet sent to 00:00:5e:00:53:01
```

## License

MIT
