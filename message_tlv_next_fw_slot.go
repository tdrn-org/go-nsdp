// message_tlv_port_statistic.go
//
// Copyright (C) 2022 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package nsdp

import (
	"fmt"
)

type NextFWSlot struct {
	Slot uint8
}

const nextFWSlotLen uint16 = 1

func EmptyNextFWSlot() *NextFWSlot {
	return NewNextFWSlot(0)
}

func NewNextFWSlot(slot uint8) *NextFWSlot {
	return &NextFWSlot{Slot: slot}
}

func unmarshalNextFWSlot(value []byte) (*NextFWSlot, error) {
	len := len(value)
	if len == 0 {
		return EmptyNextFWSlot(), nil
	}
	if len != int(nextFWSlotLen) {
		return nil, fmt.Errorf("unexpected slot length: %d", len)
	}
	return NewNextFWSlot(value[0]), nil
}

func (tlv *NextFWSlot) Type() Type {
	return TypeNextFWSlot
}

func (tlv *NextFWSlot) Length() uint16 {
	return uint16(nextFWSlotLen)
}

func (tlv *NextFWSlot) Value() []byte {
	value := make([]byte, nextFWSlotLen)
	value[0] = tlv.Slot
	return value
}

func (tlv *NextFWSlot) String() string {
	return fmt.Sprintf("NextFWSlot(%04xh) %d", TypeNextFWSlot, tlv.Slot)
}
