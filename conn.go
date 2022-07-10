// conn.go
//
// Copyright (C) 2022 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
//
package nsdp

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// Conn represents a network connection used for sending and receiving NSDP message.
type Conn struct {
	laddr              *net.UDPAddr
	taddr              *net.UDPAddr
	broadcast          bool
	host               net.HardwareAddr
	conn               *net.UDPConn
	seq                uint16
	RetryCount         int           // Number of retries during message sending
	RetryInterval      time.Duration // Interval after which message sending is retried
	ReceiveBufferSize  uint          // Receive buffer size
	ReceiveQueueLength uint          // Receive queue length
	Debug              bool
}

// NewConn establishes a new connection to the given remote endpoint.
//
// The target endpoint can be a concrete device or a broadcast address.
func NewConn(target string, broadcast bool, debug bool) (*Conn, error) {
	if debug {
		log.Printf("NSDP setting up connection...")
		log.Printf("NSDP target address: '%s' (broadcast: %t)", target, broadcast)
	}
	taddr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		return nil, err
	}
	tconn, err := net.Dial("udp", target)
	if err != nil {
		return nil, err
	}
	tconn.Close()
	lhost, _, err := net.SplitHostPort(tconn.LocalAddr().String())
	if err != nil {
		return nil, err
	}
	lport := strconv.Itoa(int(taddr.AddrPort().Port() - 1))
	listen := net.JoinHostPort(lhost, lport)
	if debug {
		log.Printf("NSDP listen address: '%s'", listen)
	}
	laddr, err := net.ResolveUDPAddr("udp", listen)
	if err != nil {
		return nil, err
	}
	host, err := lookupHardwareAddr(laddr)
	if err != nil {
		return nil, err
	}
	if debug {
		log.Printf("NSDP host MAC: '%s'", host)
	}
	lconn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, err
	}
	return &Conn{
		laddr:              laddr,
		taddr:              taddr,
		broadcast:          broadcast,
		host:               host,
		conn:               lconn,
		seq:                uint16(time.Now().Unix()),
		RetryCount:         3,
		RetryInterval:      time.Millisecond * 5000,
		ReceiveBufferSize:  8192,
		ReceiveQueueLength: 10,
		Debug:              debug,
	}, nil
}

// Close closes the connection
func (c *Conn) Close() error {
	return c.conn.Close()
}

type receiveQueueEntry struct {
	buffer []byte
	len    int
	addr   *net.UDPAddr
	err    error
}

// SendReceiveMessage sends the given NSDP message and waits for the responses.
//
// Prior to sending the message's host address and sequence number is updated according to the current connection state.
// If the broadcast flag has been set during connection creation, an arbitrary number of response message is returned.
// If the broadcast flag has not been set, exactly one response message is returned.
func (c *Conn) SendReceiveMessage(msg *Message) ([]*Message, error) {
	msg.HostAddress = c.host
	receiveQueue := make(chan receiveQueueEntry, c.ReceiveQueueLength)
	go func() {
		for {
			entry := receiveQueueEntry{
				buffer: make([]byte, c.ReceiveBufferSize),
			}
			entry.len, entry.addr, entry.err = c.conn.ReadFromUDP(entry.buffer)
			receiveQueue <- entry
			if !c.broadcast {
				break
			}
		}
	}()
	retry := 0
	retryTicker := time.NewTicker(c.RetryInterval)
	defer retryTicker.Stop()
	receivedMsgs := make([]*Message, 0)
	msg.Sequence = c.seq
	c.seq = c.seq + 1
	for retry < c.RetryCount {
		if msg.Sequence != c.seq {
			msg.Sequence = c.seq
			sendBuffer := msg.Marshal()
			if c.Debug {
				log.Printf("NSDP %s > %s:\n%s\n%s", c.laddr, c.taddr, hex.EncodeToString(sendBuffer), msg)
			}
			_, err := c.conn.WriteToUDP(sendBuffer, c.taddr)
			if err != nil {
				fmt.Println(err)
			}
		}
		select {
		case received := <-receiveQueue:
			if received.err == nil {
				receivedMsg, err := UnmarshalMessage(received.buffer[:received.len])
				if err == nil {
					if c.Debug {
						log.Printf("NSDP %s < %s:\n%s\n%s", c.laddr, received.addr, hex.EncodeToString(received.buffer[:received.len]), receivedMsg)
					}
					receivedMsgs = append(receivedMsgs, receivedMsg)
				} else {
					if c.Debug {
						log.Printf("NSDP %s < %s:\n%s", c.laddr, received.addr, hex.EncodeToString(received.buffer[:received.len]))
						log.Printf("NSDP Error while unmarshaling message; cause: %v", err)
					}
				}
			} else {
				if c.Debug {
					log.Printf("NSDP listen error: %v", received.err)
				}
			}
		case <-retryTicker.C:
			if len(receivedMsgs) != 0 {
				retry = c.RetryCount
			} else {
				c.seq = c.seq + 1
				retry = retry + 1
			}
		}
	}
	return receivedMsgs, nil
}

func lookupHardwareAddr(addr *net.UDPAddr) (net.HardwareAddr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		ifaceAddrs, err := iface.Addrs()
		if err == nil {
			for _, ifaceAddr := range ifaceAddrs {
				if strings.HasPrefix(ifaceAddr.String(), addr.IP.String()) {
					return iface.HardwareAddr, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("failed to lookup hardware address for address: %s", addr)
}
