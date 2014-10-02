//package main //for test
package hpack

import (
	"encoding/hex"
	"fmt"
	"huffman"
)

func PackIntRepresentation(I uint32, N byte) (buf *[]byte) {
	if I < uint32(1<<N)-1 {
		return &[]byte{byte(I)}
	} else {
		buf = &[]byte{byte(1<<N) - 1}
		I -= uint32(1<<N) - 1
		for I >= 0x80 {
			*buf = append(*buf, byte(I)&0x7f|0x80)
			I = (I >> 7)
		}
		*buf = append(*buf, byte(I))
		return buf
	}
}

func ParseIntRepresentation(buf []byte, N byte) (I, cursor uint32) {
	I = uint32(buf[0] & ((1 << N) - 1)) // byte could be used as byte
	cursor = 1
	if I < ((1 << N) - 1) {
		return I, cursor
	} else {
		var M byte = 0
		for (buf[cursor] & 0x80) > 0 {
			I += uint32(buf[cursor]&0x7f) * (1 << M)
			M += 7
			cursor += 1
		}
		I += uint32(buf[cursor]&0x7f) * (1 << M)
		return I, cursor + 1
	}
}

func ExtractContent(buf []byte, length uint32, isHuffman bool) string {
	if isHuffman {
		return huffman.Root.Decode(buf, length)
	} else {
		return string(buf[:length])
	}
}

func ParseFromByte(buf []byte) (content string, cursor uint32) {
	isHuffman := false
	if buf[0]&0x80 > 0 {
		isHuffman = true
	}
	length, cursor := ParseIntRepresentation(buf, 7)
	content = ExtractContent(buf[cursor:], length, isHuffman)
	cursor += length
	return
}

func ParseHeader(index uint32, buf []byte, isIndexed bool) (name, value string, cursor uint32) {
	if c := uint32(0); !isIndexed {
		if index == 0 {
			name, c = ParseFromByte(buf[cursor:])
			cursor += c
		}
		value, c = ParseFromByte(buf[cursor:])
		cursor += c
	}

	if index > 0 {
		header := GetHeader(index)

		name = header.Name
		if len(value) == 0 {
			value = header.Value
		}
	}
	return
}

func Decode(wire string) (Headers []Header) {
	var buf *[]byte
	nums, err := hex.DecodeString(string(wire))
	if err != nil {
		panic(err)
	}
	buf = &nums

	var cursor uint32 = 0
	for cursor < uint32(len(nums)) {
		isIndexed := false
		isIncremental := false
		var index, c uint32
		if (*buf)[cursor]&0xe0 == 0x20 {
			// 7.3 Header Table Size Update
			size, _ := ParseIntRepresentation((*buf)[cursor:], 5)
			SetMaxHeaderTableSize(size)
			cursor += 1
		}

		if ((*buf)[cursor] & 0x80) > 0 {
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

		name, value, c := ParseHeader(index, (*buf)[cursor:], isIndexed)
		cursor += c

		header := Header{name, value}
		if isIncremental {
			AddHeader(header)
		}
		Headers = append(Headers, header)
	}

	return
}

func main() {
	//nums, _ := hex.DecodeString(string("1FA18DB701"))
	//fmt.Println(nums)
	//fmt.Println(ParseIntRepresentation(nums, 5))
	//decode("ff80000111")
	fmt.Println(Decode("00073a6d6574686f640347455400073a736368656d650468747470000a3a617574686f726974790f7777772e7961686f6f2e636f2e6a7000053a70617468012f"))
	fmt.Println(huffman.HUFFMAN_TABLE)
}
