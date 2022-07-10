// conn_test.go
//
// Copyright (C) 2022 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
//
package nsdp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnSendReceiveMessage(t *testing.T) {
	conn, err := NewConn("255.255.255.255:63322", true, true)
	require.Nil(t, err)
	require.Nil(t, err)
	message := NewMessage(ReadRequest)
	message.AppendTLV(EmptyDeviceModel())
	message.AppendTLV(EmptyDeviceName())
	message.AppendTLV(EmptyDeviceMAC())
	message.AppendTLV(EmptyDeviceLocation())
	message.AppendTLV(EmptyDeviceIP())
	message.AppendTLV(EmptyDeviceNetmask())
	message.AppendTLV(EmptyRouterIP())
	message.AppendTLV(EmptyPortStatus())
	message.AppendTLV(EmptyPortStatistic())
	conn.SendReceiveMessage(message)
}
