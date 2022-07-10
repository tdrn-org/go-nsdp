// message_tlv_port_statistic.go
//
// Copyright (C) 2022 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
//
package nsdp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PortStatistic struct {
	Port       uint8
	Received   uint64
	Send       uint64
	Packets    uint64
	Broadcasts uint64
	Multicasts uint64
	Errors     uint64
}

const portStatisticLen uint16 = 49

func EmptyPortStatistic() *PortStatistic {
	return &PortStatistic{}
}

func NewPortStatistic(port uint8, received uint64, send uint64, packets uint64, broadcasts uint64, multicasts uint64, errors uint64) *PortStatistic {
	return &PortStatistic{
		Port:       port,
		Received:   received,
		Send:       send,
		Packets:    packets,
		Broadcasts: broadcasts,
		Multicasts: multicasts,
		Errors:     errors,
	}
}

func unmarshalPortStatistic(value []byte) (*PortStatistic, error) {
	len := len(value)
	if len != int(portStatisticLen) {
		return nil, fmt.Errorf("unexpected port statistic length: %d", len)
	}
	buffer := bytes.NewBuffer(value)
	tlv := EmptyPortStatistic()
	tlv.Port, _ = buffer.ReadByte()
	binary.Read(buffer, binary.BigEndian, &tlv.Received)
	binary.Read(buffer, binary.BigEndian, &tlv.Send)
	binary.Read(buffer, binary.BigEndian, &tlv.Packets)
	binary.Read(buffer, binary.BigEndian, &tlv.Broadcasts)
	binary.Read(buffer, binary.BigEndian, &tlv.Multicasts)
	binary.Read(buffer, binary.BigEndian, &tlv.Errors)
	return tlv, nil
}

func (tlv *PortStatistic) Type() Type {
	return TypePortStatistic
}

func (tlv *PortStatistic) Length() uint16 {
	return uint16(portStatisticLen)
}

func (tlv *PortStatistic) Value() []byte {
	value := make([]byte, portStatisticLen)
	buffer := bytes.NewBuffer(value)
	buffer.WriteByte(tlv.Port)
	binary.Write(buffer, binary.BigEndian, tlv.Received)
	binary.Write(buffer, binary.BigEndian, tlv.Send)
	binary.Write(buffer, binary.BigEndian, tlv.Packets)
	binary.Write(buffer, binary.BigEndian, tlv.Broadcasts)
	binary.Write(buffer, binary.BigEndian, tlv.Multicasts)
	binary.Write(buffer, binary.BigEndian, tlv.Errors)
	return value
}

func (tlv *PortStatistic) String() string {
	return fmt.Sprintf("PortStatistic(%04xh) Port%d Received: %d, Send: %d, Packets: %d, Broadcasts: %d, Multicasts: %d, Errors: %d", TypePortStatistic, tlv.Port, tlv.Received, tlv.Send, tlv.Packets, tlv.Broadcasts, tlv.Multicasts, tlv.Errors)
}
