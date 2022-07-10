// message_tlv.go
//
// Copyright (C) 2022 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
//
package nsdp

import "fmt"

type Type uint16

// TLV message element types
const (
	TypeDeviceModel    Type = 0x0001
	TypeDeviceName     Type = 0x0003
	TypeDeviceMAC      Type = 0x0004
	TypeDeviceLocation Type = 0x0005
	TypeDeviceIP       Type = 0x0006
	TypeDeviceNetmask  Type = 0x0007
	TypeRouterIP       Type = 0x0008
	TypePassword       Type = 0x000a
	TypeDHCPMode       Type = 0x000b
	TypeFWVersionSlot1 Type = 0x000d
	TypeFWVersionSlot2 Type = 0x000e
	TypeNextFWSlot     Type = 0x000f
	TypePortStatus     Type = 0x0c00
	TypePortStatistic  Type = 0x1000
	TypeGetVlanInfo    Type = 0x2800
	TypeDeleteVlan     Type = 0x2c00
	TypeEOM            Type = 0xffff // EOM marker prefix (always the last TLV and automatically part of each message)
)

// Type-length-value message element's interface
type TLV interface {
	Type() Type
	Length() uint16
	Value() []byte
}

func unmarshalTLV(tlvType uint16, tlvValue []byte) (TLV, error) {
	switch tlvType {
	case uint16(TypeDeviceModel):
		return unmarshalDeviceModel(tlvValue)
	case uint16(TypeDeviceName):
		return unmarshalDeviceName(tlvValue)
	case uint16(TypeDeviceMAC):
		return unmarshalDeviceMAC(tlvValue)
	case uint16(TypeDeviceLocation):
		return unmarshalDeviceLocation(tlvValue)
	case uint16(TypeDeviceIP):
		return unmarshalDeviceIP(tlvValue)
	case uint16(TypeDeviceNetmask):
		return unmarshalDeviceNetmask(tlvValue)
	case uint16(TypeRouterIP):
		return unmarshalRouterIP(tlvValue)
	case uint16(TypePortStatus):
		return unmarshalPortStatus(tlvValue)
	case uint16(TypePortStatistic):
		return unmarshalPortStatistic(tlvValue)
	}
	return nil, fmt.Errorf("unrecognized TLV type: %04xh", tlvType)
}
