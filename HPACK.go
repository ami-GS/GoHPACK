package main //for test

import (
	"encoding/hex"
	"fmt"
)

type Header struct {
	Name, Value string
}

func ParseIntRepresentation(buf []byte, N byte) (I int, cursor byte) {
	I = int(buf[0] & ((1 << N) - 1)) // byte could be used as byte
	cursor = 1
	if I < ((1 << N) - 1) {
		return I, cursor
	} else {
		var M byte = 0
		for (buf[cursor] & 0x80) > 0 {
			I += int(buf[cursor]&0x7f) * (1 << M)
			M += 7
			cursor += 1
		}
		I += int(buf[cursor]&0x7f) * (1 << M)
		return I, cursor + 1
	}
}

func ExtractContent(buf []byte, length int, isHuffman bool) (content string) {
	if isHuffman {
		return content
	} else {
		content = string(buf[:length])
		return content
	}
}

func ParseFromByte(buf []byte) (content string, cursor byte) {
	isHuffman := false
	if buf[0]&0x80 > 0 {
		isHuffman = true
	}
	length, cursor := ParseIntRepresentation(buf, 7)
	content = ExtractContent(buf[cursor:], length, isHuffman)
	cursor += byte(length)
	return content, cursor
}

func ParseHeader(index int, table int, buf []byte, isIndexed bool) (name, value string, cursor byte) {
	if c := byte(0); !isIndexed {
		if index == 0 {
			name, c = ParseFromByte(buf[cursor:])
			cursor += c
		}
		value, c = ParseFromByte(buf[cursor:])
		cursor += c
	}

	if index > 0 {
		//get header from table
	}

	return name, value, cursor
}

func decode(wire string) []Header {
	var Headers []Header
	var buf *[]byte
	nums, err := hex.DecodeString(string(wire))
	if err != nil {
		panic(err)
	}
	buf = &nums

	var cursor byte = 0
	for cursor < byte(len(nums)) {
		isIndexed := false
		isIncremental := false
		var index int
		var c byte
		if (*buf)[cursor]&0xe0 == 0x20 {
			// 7.3 Header Table Size Update
			cursor = 1
		} else if ((*buf)[cursor] & 0x80) > 0 {
			// 7.1 Indexed Header Field
			if ((*buf)[cursor] & 0x7f) == 0 {
				panic('a')
			}
			index, c = ParseIntRepresentation((*buf)[cursor:], 7)
			isIndexed = true
		} else {
			if (*buf)[cursor]&0xc0 == 0x40 {
				// 7.2.1 Literal Header Field with Incremental Indexing
				index, c = ParseIntRepresentation((*buf)[cursor:], 6)
				isIncremental = true
			} else if (*buf)[cursor]&0xf0 == 0xf0 {
				// 7.2.3 Literal Header Field never Indexed
				index, c = ParseIntRepresentation((*buf)[cursor:], 4)
			} else {
				// 7.2.2 Literal Header Field without Indexing
				index, c = ParseIntRepresentation((*buf)[cursor:], 4)
			}
		}
		cursor += c

		table := 1 //for test
		name, value, c := ParseHeader(index, table, (*buf)[cursor:], isIndexed)
		cursor += c

		if isIncremental {
			//add to table
		}
		Headers = append(Headers, []Header{{name, value}}...)
	}

	//d := hex.EncodeToString(nums)
	//fmt.Println(d)
	return Headers
}

func main() {
	//nums, _ := hex.DecodeString(string("1FA18DB701"))
	//fmt.Println(nums)
	//fmt.Println(ParseIntRepresentation(nums, 5))
	//decode("ff80000111")
	fmt.Println(decode("00073a6d6574686f640347455400073a736368656d650468747470000a3a617574686f726974790f7777772e7961686f6f2e636f2e6a7000053a70617468012f"))
	fmt.Println()
}
