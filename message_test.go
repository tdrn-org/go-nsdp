// message_test.go
//
// Copyright (C) 2022 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
//
package nsdp

import (
	"crypto/rand"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeviceModelMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, NewDeviceModel("Model"))
}

func TestDeviceModelString(t *testing.T) {
	runMessageStringTest(t, NewDeviceModel("Model"), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: DeviceModel(0001h) 'Model'\nEOM   : ffff0000h")
}

func TestDeviceNameMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, NewDeviceName("Name"))
}

func TestDeviceNameString(t *testing.T) {
	runMessageStringTest(t, NewDeviceName("Name"), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: DeviceName(0003h) 'Name'\nEOM   : ffff0000h")
}

func TestDeviceMAClMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, NewDeviceMAC(getRandomMAC()))
}

func TestDeviceMACString(t *testing.T) {
	runMessageStringTest(t, NewDeviceMAC(getStaticMAC()), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: DeviceMAC(0004h) 01:02:03:04:05:06\nEOM   : ffff0000h")
}

func TestDeviceLocationlMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, NewDeviceLocation("DeviceLocation"))
}

func TestDeviceLocationString(t *testing.T) {
	runMessageStringTest(t, NewDeviceLocation("Location"), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: DeviceLocation(0005h) 'Location'\nEOM   : ffff0000h")
}

func TestDeviceIPMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, NewDeviceIP(getRandomIP()))
}

func TestDeviceIPString(t *testing.T) {
	runMessageStringTest(t, NewDeviceIP(getStaticIP()), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: DeviceIP(0006h) 1.2.3.4\nEOM   : ffff0000h")
}

func TestDeviceNetmaskMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, NewDeviceNetmask(getRandomIP()))
}

func TestDeviceNetmaskString(t *testing.T) {
	runMessageStringTest(t, NewDeviceNetmask(getStaticIP()), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: DeviceNetmask(0007h) 1.2.3.4\nEOM   : ffff0000h")
}

func TestRouterIPMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, NewRouterIP(getRandomIP()))
}

func TestRouterIPString(t *testing.T) {
	runMessageStringTest(t, NewRouterIP(getStaticIP()), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: RouterIP(0008h) 1.2.3.4\nEOM   : ffff0000h")
}

func TestPortStatusMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, NewPortStatus(1, 2))
}

func TestPortStatusString(t *testing.T) {
	runMessageStringTest(t, NewPortStatus(1, 2), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: PortStatus(0c00h) Port1 Status: 10Mbit/full-duplex Unknown1: 00h\nEOM   : ffff0000h")
}

func TestPortStatisticMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, NewPortStatistic(1, 2, 3, 4, 5, 6, 7))
}

func TestPortStatisticString(t *testing.T) {
	runMessageStringTest(t, NewPortStatistic(1, 2, 3, 4, 5, 6, 7), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: PortStatistic(1000h) Port1 Received: 2, Send: 3, Packets: 4, Broadcasts: 5, Multicasts: 6, Errors: 7\nEOM   : ffff0000h")
}

func runMessageMarshalingTest(t *testing.T, tlv TLV) {
	message1 := NewMessage(ReadResponse)
	message1.AppendTLV(tlv)
	marshaledBytes := message1.Marshal()
	message2, err := UnmarshalMessage(marshaledBytes)
	require.Nil(t, err)
	unmarshaledBytes := message2.Marshal()
	require.Equal(t, marshaledBytes, unmarshaledBytes)
}

func runMessageStringTest(t *testing.T, tlv TLV, expected string) {
	message := NewMessage(ReadResponse)
	message.AppendTLV(tlv)
	messageString := fmt.Sprint(message)
	require.Equal(t, expected, messageString)
}

func getRandomMAC() net.HardwareAddr {
	mac := make([]byte, 6)
	rand.Read(mac)
	return mac
}

func getStaticMAC() net.HardwareAddr {
	var mac = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	return mac
}

func getRandomIP() net.IP {
	ip := make([]byte, 4)
	rand.Read(ip)
	return ip
}

func getStaticIP() net.IP {
	var ip = []byte{0x01, 0x02, 0x03, 0x04}
	return ip
}
