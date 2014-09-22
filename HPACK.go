package main //for test

import (
	"encoding/hex"
	"fmt"
)

type Header struct {
	Name, Value string
}

func ParseIntRepresentation(buf []byte, N byte) (index byte, c int) {
	I := (buf[0] & ((1 << N) - 1))
	cursor := 1
	if I < ((1 << N) - 1) {
		return I, cursor
	} else {
		var M byte = 0
		for (buf[cursor] & 0x80) > 0 {
			I += (buf[cursor] & 0x7f) * (1 << M)
			M += 7
			cursor += 1
		}
		I += (buf[cursor] & 0x7f) * (1 << M)
		return I, cursor + 1
	}
}

func ParseHeader() {

}

func decode(wire string) []Header {
	var Headers []Header
	var buf *[]byte
	nums, err := hex.DecodeString(string(wire))
	if err != nil {
		panic(err)
	}
	buf = &nums

	//for i, v := range *buf {
	//	fmt.Println(i, v)
	//}

	cursor := 0
	for cursor < len(nums) {
		isIndexed := false
		isIncremental := false
		var index byte
		if (*buf)[cursor]&0xe0 == 0x20 {
			//change table size
		} else if ((*buf)[cursor] & 0x80) > 0 {
			if ((*buf)[cursor] & 0x7f) == 0 {
				panic('a')
			}
			index, c := ParseIntRepresentation((*buf)[cursor:], 7)
			cursor += c
			isIndexed = true
		} else {
			c := 0
			if (*buf)[cursor]&0xc0 == 0x40 {
				index, c := ParseIntRepresentation((*buf)[cursor:], 6)
				isIncremental = true
			} else if (*buf)[cursor]&0xf0 == 0xf0 {
				index, c := ParseIntRepresentation((*buf)[cursor:], 4)
			} else {
				index, c := ParseIntRepresentation((*buf)[cursor:], 4)
			}
			cursor += c
		}
		table := 1 //for test
		name, value, c := ParseHeader(index, table, (*buf)[cursor:], isIndexed)
		cursor += c

		if isIncremental {
			//add to table
		}
		Headers.append(Header{name, value})
	}

	//d := hex.EncodeToString(nums)
	//fmt.Println(d)
	return Headers
}

func main() {
	decode("ff80000111")
	fmt.Println()
}
