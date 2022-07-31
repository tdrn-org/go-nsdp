// message.go
//
// Copyright (C) 2022 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package nsdp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

// NSDP message (see https://en.wikipedia.org/wiki/Netgear_Switch_Discovery_Protocol).
//
// A message is constructed from a Header an EOM marker as well as an arbitary number of TLV (type-length-value) payload elements.
// The Header defines the general message processing rules (espcially type of operation and target device). The TLV elements define
// the actual message content.
type Message struct {
	Header *Header // Message header
	Body   []TLV   // Message body (payload)
	EOM    *EOM    // End-of-message marker
}

// NewMessage constructs a new message for the given operation code with an empty list of TLVs.
func NewMessage(operation OperationCode) *Message {
	return &Message{
		Header: newHeader(operation),
		Body:   make([]TLV, 0),
		EOM:    newEOM(),
	}
}

func (m *Message) prepareMessage(hostAddress net.HardwareAddr, sequence Sequence) *Message {
	return &Message{
		Header: m.Header.prepareHeader(hostAddress, sequence),
		Body:   m.Body,
		EOM:    m.EOM,
	}
}

// AppendTLV updates the message by appending an additional TLV to it.
func (m *Message) AppendTLV(tlv TLV) {
	m.Body = append(m.Body, tlv)
}

func (m *Message) String() string {
	builder := &strings.Builder{}
	m.Header.writeString(builder)
	builder.WriteRune('\n')
	for i, tlv := range m.Body {
		builder.WriteString(fmt.Sprintf("TLV[%d]: %s\n", i, tlv))
	}
	m.EOM.writeString(builder)
	return builder.String()
}

// Marshal encodes the message to its NSDP compliant byte stream.
func (m *Message) Marshal() []byte {
	buffer := &bytes.Buffer{}
	m.MarshalBuffer(buffer)
	return buffer.Bytes()
}

// MarshalBuffer encodes the message to its NSDP compliant byte stream.
func (m *Message) MarshalBuffer(buffer *bytes.Buffer) {
	m.Header.marshalBuffer(buffer)
	for _, tlv := range m.Body {
		binary.Write(buffer, binary.BigEndian, tlv.Type())
		if m.Header.Operation == ReadRequest {
			binary.Write(buffer, binary.BigEndian, uint16(0))
		} else {
			binary.Write(buffer, binary.BigEndian, uint16(tlv.Length()))
			buffer.Write(tlv.Value())
		}
	}
	m.EOM.marshalBuffer(buffer)
}

// UnmarshalMessage decodes a message from the given NSDP byte stream.
func UnmarshalMessage(buf []byte) (*Message, error) {
	buffer := bytes.NewBuffer(buf)
	return UnmarshalMessageBuffer(buffer)
}

// UnmarshalMessage decodes a message from the given NSDP byte stream.
func UnmarshalMessageBuffer(buffer *bytes.Buffer) (*Message, error) {
	header, err := unmarshalHeaderBuffer(buffer)
	if err != nil {
		return nil, err
	}
	tlvs := make([]TLV, 0)
	for {
		var tlvType uint16
		err = binary.Read(buffer, binary.BigEndian, &tlvType)
		if err != nil {
			return nil, fmt.Errorf("error while decoding TLV type; cause: %v", err)
		}
		var tlvLength uint16
		err = binary.Read(buffer, binary.BigEndian, &tlvLength)
		if err != nil {
			return nil, fmt.Errorf("error while decoding TLV type; cause: %v", err)
		}
		if tlvType == uint16(TypeEOM) {
			if tlvLength != 0 {
				return nil, fmt.Errorf("unexpected EOM marker: %04x%04xh", tlvType, tlvLength)
			}
			break
		}
		tlvValue := make([]byte, tlvLength)
		_, err = buffer.Read(tlvValue)
		if err != nil {
			return nil, fmt.Errorf("error while decoding TLV value; cause: %v", err)
		}
		tlv, err := unmarshalTLV(tlvType, tlvValue)
		if err != nil {
			return nil, fmt.Errorf("error while decoding TLV type %04xh; cause: %v", tlvType, err)
		}
		tlvs = append(tlvs, tlv)
	}
	return &Message{
		Header: header,
		Body:   tlvs,
		EOM:    newEOM(),
	}, nil
}
