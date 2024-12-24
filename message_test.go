// message_test.go
//
// Copyright (C) 2022-2024 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package nsdp_test

import (
	"crypto/rand"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-nsdp"
)

func TestDeviceModelMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, nsdp.NewDeviceModel("Model"))
}

func TestDeviceModelString(t *testing.T) {
	runMessageStringTest(t, nsdp.NewDeviceModel("Model"), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: DeviceModel(0001h) 'Model'\nEOM   : ffff0000h")
}

func TestDeviceNameMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, nsdp.NewDeviceName("Name"))
}

func TestDeviceNameString(t *testing.T) {
	runMessageStringTest(t, nsdp.NewDeviceName("Name"), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: DeviceName(0003h) 'Name'\nEOM   : ffff0000h")
}

func TestDeviceMAClMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, nsdp.NewDeviceMAC(getRandomMAC()))
}

func TestDeviceMACString(t *testing.T) {
	runMessageStringTest(t, nsdp.NewDeviceMAC(getStaticMAC()), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: DeviceMAC(0004h) 01:02:03:04:05:06\nEOM   : ffff0000h")
}

func TestDeviceLocationlMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, nsdp.NewDeviceLocation("DeviceLocation"))
}

func TestDeviceLocationString(t *testing.T) {
	runMessageStringTest(t, nsdp.NewDeviceLocation("Location"), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: DeviceLocation(0005h) 'Location'\nEOM   : ffff0000h")
}

func TestDeviceIPMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, nsdp.NewDeviceIP(getRandomIP()))
}

func TestDeviceIPString(t *testing.T) {
	runMessageStringTest(t, nsdp.NewDeviceIP(getStaticIP()), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: DeviceIP(0006h) 1.2.3.4\nEOM   : ffff0000h")
}

func TestDeviceNetmaskMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, nsdp.NewDeviceNetmask(getRandomIP()))
}

func TestDeviceNetmaskString(t *testing.T) {
	runMessageStringTest(t, nsdp.NewDeviceNetmask(getStaticIP()), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: DeviceNetmask(0007h) 1.2.3.4\nEOM   : ffff0000h")
}

func TestRouterIPMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, nsdp.NewRouterIP(getRandomIP()))
}

func TestRouterIPString(t *testing.T) {
	runMessageStringTest(t, nsdp.NewRouterIP(getStaticIP()), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: RouterIP(0008h) 1.2.3.4\nEOM   : ffff0000h")
}

func TestDHCPModeMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, nsdp.NewDHCPMode(1))
}

func TestDHCPModeString(t *testing.T) {
	runMessageStringTest(t, nsdp.NewDHCPMode(1), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: DHCPMode(000bh) Enabled\nEOM   : ffff0000h")
}

func TestFWVersionSlot1Marshaling(t *testing.T) {
	runMessageMarshalingTest(t, nsdp.NewFWVersionSlot1("1.2.3.4"))
}

func TestFWVersionSlot1String(t *testing.T) {
	runMessageStringTest(t, nsdp.NewFWVersionSlot1("1.2.3.4"), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: FWVersionSlot1(000dh) '1.2.3.4'\nEOM   : ffff0000h")
}

func TestFWVersionSlot2Marshaling(t *testing.T) {
	runMessageMarshalingTest(t, nsdp.NewFWVersionSlot2("4.3.2.1"))
}

func TestFWVersionSlot2String(t *testing.T) {
	runMessageStringTest(t, nsdp.NewFWVersionSlot2("4.3.2.1"), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: FWVersionSlot2(000eh) '4.3.2.1'\nEOM   : ffff0000h")
}

func TestPortStatusMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, nsdp.NewPortStatus(1, 2))
}

func TestPortStatusString(t *testing.T) {
	runMessageStringTest(t, nsdp.NewPortStatus(1, 2), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: PortStatus(0c00h) Port1 Status: 10Mbit/full-duplex Unknown1: 00h\nEOM   : ffff0000h")
}

func TestPortStatisticMarshaling(t *testing.T) {
	runMessageMarshalingTest(t, nsdp.NewPortStatistic(1, 2, 3, 4, 5, 6, 7))
}

func TestPortStatisticString(t *testing.T) {
	runMessageStringTest(t, nsdp.NewPortStatistic(1, 2, 3, 4, 5, 6, 7), "Header: 01h 02h 0000h 00000000h 00:00:00:00:00:00 00:00:00:00:00:00 0000h 0000h 4e534450h\nTLV[0]: PortStatistic(1000h) Port1 Received: 2, Sent: 3, Packets: 4, Broadcasts: 5, Multicasts: 6, Errors: 7\nEOM   : ffff0000h")
}

func runMessageMarshalingTest(t *testing.T, tlv nsdp.TLV) {
	runRequestMessageMarshalingTest(t, tlv)
	runResponseMessageMarshalingTest(t, tlv)
}

func runRequestMessageMarshalingTest(t *testing.T, tlv nsdp.TLV) {
	message1 := nsdp.NewMessage(nsdp.ReadRequest)
	message1.AppendTLV(tlv)
	marshaledBytes := message1.Marshal()
	message2, err := nsdp.UnmarshalMessage(marshaledBytes)
	require.NoError(t, err)
	unmarshaledBytes := message2.Marshal()
	require.Equal(t, marshaledBytes, unmarshaledBytes)
}
func runResponseMessageMarshalingTest(t *testing.T, tlv nsdp.TLV) {
	message1 := nsdp.NewMessage(nsdp.ReadResponse)
	message1.AppendTLV(tlv)
	marshaledBytes := message1.Marshal()
	message2, err := nsdp.UnmarshalMessage(marshaledBytes)
	require.NoError(t, err)
	unmarshaledBytes := message2.Marshal()
	require.Equal(t, marshaledBytes, unmarshaledBytes)
}

func runMessageStringTest(t *testing.T, tlv nsdp.TLV, expected string) {
	message := nsdp.NewMessage(nsdp.ReadResponse)
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
