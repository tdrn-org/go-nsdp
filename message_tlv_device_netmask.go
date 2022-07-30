// message_tlv_device_netmask.go
//
// Copyright (C) 2022 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package nsdp

import (
	"fmt"
	"net"
)

type DeviceNetmask struct {
	Netmask net.IP
}

func EmptyDeviceNetmask() *DeviceNetmask {
	return NewDeviceNetmask(net.IP{})
}

func NewDeviceNetmask(netmask net.IP) *DeviceNetmask {
	return &DeviceNetmask{Netmask: netmask}
}

func unmarshalDeviceNetmask(bytes []byte) (*DeviceNetmask, error) {
	len := len(bytes)
	if len == 0 {
		return EmptyDeviceNetmask(), nil
	}
	if len != 4 && len != 16 {
		return nil, fmt.Errorf("unexpected device netmask length: %d", len)
	}
	return NewDeviceNetmask(net.IP(bytes)), nil
}

func (tlv *DeviceNetmask) Type() Type {
	return TypeDeviceNetmask
}

func (tlv *DeviceNetmask) Length() uint16 {
	return uint16(len(tlv.Netmask))
}

func (tlv *DeviceNetmask) Value() []byte {
	return tlv.Netmask
}

func (tlv *DeviceNetmask) String() string {
	return fmt.Sprintf("DeviceNetmask(%04xh) %s", TypeDeviceNetmask, tlv.Netmask)
}
