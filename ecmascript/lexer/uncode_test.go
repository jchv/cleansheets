package lexer

import (
	"reflect"
	"strconv"
	"testing"
)

func TestEncodeUTF16(t *testing.T) {
	tests := []struct {
		s string
		u []uint16
	}{
		{"", []uint16{}},
		{" ", []uint16{0x0020}},
		{"\u0000", []uint16{0x0000}},
		{"\u0080", []uint16{0x0080}},
		{"\u0800", []uint16{0x0800}},
		{"\U00010001", []uint16{0xd800, 0xdc01}},
		{"\U00018888", []uint16{0xd822, 0xdc88}},
		{"\U0001aaaa", []uint16{0xd82a, 0xdeaa}},
		{"\U0001ffff", []uint16{0xd83f, 0xdfff}},
		{"\U00020000", []uint16{0xd840, 0xdc00}},
		{"\U0010ffff", []uint16{0xdbff, 0xdfff}},
		{"test", []uint16{0x74, 0x65, 0x73, 0x74}},
		{"日本語", []uint16{0x65e5, 0x672c, 0x8a9e}},
		{"💌 ✉️", []uint16{0xd83d, 0xdc8c, 0x20, 0x2709, 0xfe0f}},
	}

	for _, test := range tests {
		t.Run(strconv.Quote(test.s), func(t *testing.T) {
			result := EncodeUTF16(test.s)
			if !reflect.DeepEqual(result, test.u) {
				t.Errorf("EncodeUTF16(%q) = %v != %v", test.s, result, test.u)
			}
		})
	}
}
