// conn.go
//
// Copyright (C) 2022 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package nsdp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// Conn represents a network connection used for sending and receiving NSDP messages.
type Conn struct {
	laddr              *net.UDPAddr
	taddr              *net.UDPAddr
	host               net.HardwareAddr
	conn               *net.UDPConn
	seq                Sequence
	ReceiveBufferSize  uint // Receive buffer size
	ReceiveQueueLength uint // Receive queue length
	ReceiveTimeout     time.Duration
	Debug              bool // Enables debug output via log.Printf
}

// NewConn establishes a new connection to the given remote target.
func NewConn(target string, debug bool) (*Conn, error) {
	if debug {
		log.Printf("NSDP setting up connection...")
		log.Printf("NSDP target address: '%s'", target)
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
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, err
	}
	return &Conn{
		laddr:              laddr,
		taddr:              taddr,
		host:               host,
		conn:               conn,
		seq:                Sequence(time.Now().UnixNano()),
		ReceiveBufferSize:  8192,
		ReceiveQueueLength: 16,
		ReceiveTimeout:     time.Millisecond * 2000,
		Debug:              debug,
	}, nil
}

// Close closes the connection.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// SendReceiveMessage sends the given NSDP message and waits for responses.
//
// Prior to sending the message's host address and sequence number is updated according to the current connection state.
// If the message's device address is empty, an arbitrary number of response message is returned.
// If the message's device address has been set, exactly one response message is returned.
func (c *Conn) SendReceiveMessage(msg *Message) (map[string]*Message, error) {
	c.seq += 1
	c.conn.SetReadDeadline(time.Now().Add(c.ReceiveTimeout))
	if bytes.Equal(msg.Header.DeviceAddress, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) {
		return c.sendReceiveBroadcastMessage(msg)
	}
	return c.sendReceiveUnicastMessage(msg)
}

type receiveQueueEntry struct {
	msg *Message
	err error
}

func (c *Conn) sendReceiveBroadcastMessage(msg *Message) (map[string]*Message, error) {
	receiveQueue := make(chan *receiveQueueEntry, c.ReceiveQueueLength)
	go func() {
		for {
			msg, err := c.receiveMessage()
			receiveQueue <- &receiveQueueEntry{
				msg: msg,
				err: err,
			}
			if err != nil {
				break
			}
		}
	}()
	err := c.sendMessage(msg)
	if err != nil {
		return nil, err
	}
	receivedMsgs := make(map[string]*Message, 0)
	for {
		received := <-receiveQueue
		if received.err != nil {
			if nerr, ok := received.err.(net.Error); ok && nerr.Timeout() {
				break
			}
			return nil, received.err
		}
		receivedMsgs[received.msg.Header.DeviceAddress.String()] = received.msg
	}
	return receivedMsgs, nil
}

func (c *Conn) sendReceiveUnicastMessage(msg *Message) (map[string]*Message, error) {
	receiveQueue := make(chan *receiveQueueEntry, 1)
	go func() {
		for {
			msg, err := c.receiveMessage()
			receiveQueue <- &receiveQueueEntry{
				msg: msg,
				err: err,
			}
			break
		}
	}()
	err := c.sendMessage(msg)
	if err != nil {
		return nil, err
	}
	receivedMsgs := make(map[string]*Message, 0)
	received := <-receiveQueue
	if received.err != nil {
		return nil, received.err
	}
	receivedMsgs[received.msg.Header.DeviceAddress.String()] = received.msg
	return receivedMsgs, nil
}

func (c *Conn) sendMessage(msg *Message) error {
	preparedMsg := msg.prepareMessage(c.host, c.seq)
	sendBuffer := preparedMsg.Marshal()
	if c.Debug {
		log.Printf("NSDP %s > %s:\n%s\n%s", c.laddr, c.taddr, hex.EncodeToString(sendBuffer), preparedMsg)
	}
	_, err := c.conn.WriteToUDP(sendBuffer, c.taddr)
	return err
}

func (c *Conn) receiveMessage() (*Message, error) {
	buffer := make([]byte, c.ReceiveBufferSize)
	for {
		len, addr, err := c.conn.ReadFromUDP(buffer)
		if err != nil {
			return nil, err
		}
		msg, err := UnmarshalMessage(buffer[:len])
		if err != nil {
			if c.Debug {
				log.Printf("NSDP %s < %s:\n%s", c.laddr, addr, hex.EncodeToString(buffer[:len]))
				log.Printf("NSDP Error while unmarshaling message; cause: %v", err)
			}
			return nil, err
		}
		if msg.Header.Sequence == c.seq {
			if c.Debug {
				log.Printf("NSDP %s < %s:\n%s\n%s", c.laddr, addr, hex.EncodeToString(buffer[:len]), msg)
			}
			return msg, nil
		} else {
			if c.Debug {
				log.Printf("NSDP %s < %s:\nIgnoring unsolicted message (sequence: %04xh)", c.laddr, addr, msg.Header.Sequence)
			}
		}
	}
}

func lookupHardwareAddr(addr *net.UDPAddr) (net.HardwareAddr, error) {
	// lo has no real MAC; use 00:00:00:00:00:00
	if addr.IP.IsLoopback() {
		return make([]byte, 6), nil
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		ifaceAddrs, err := iface.Addrs()
		if err == nil {
			for _, ifaceAddr := range ifaceAddrs {
				if strings.HasPrefix(ifaceAddr.String(), addr.IP.String()) {
					if len(iface.HardwareAddr) != 6 {
						return nil, fmt.Errorf("failed to lookup hardware address for interface %s", iface.Name)
					}
					return iface.HardwareAddr, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("failed to lookup hardware address for address: %s", addr)
}
