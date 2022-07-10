// message_tlv.go
//
// Copyright (C) 2022 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
//
package nsdp

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
