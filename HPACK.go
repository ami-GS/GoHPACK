//package main //for test
package GoHPACK

import (
	"encoding/hex"
	"fmt"
	"github.com/ami-GS/GoHPACK/huffman"
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

func PackContent(content string, toHuffman bool) string {
	if len(content) == 0 {
		if toHuffman {
			return "80"
		} else {
			return "00"
		}
	}

	Wire := ""
	if toHuffman {
		encoded, length := huffman.Root.Encode(content)
		intRep := PackIntRepresentation(uint32(length), 7)
		(*intRep)[0] |= 0x80

		//Wire += hex.EncodeToString(*intRep) + strings.Trim(hex.EncodeToString(b), "00") // + encoded
		Wire += hex.EncodeToString(*intRep) + hex.EncodeToString(encoded)
	} else {
		intRep := PackIntRepresentation(uint32(len(content)), 7)
		Wire += hex.EncodeToString(*intRep) + hex.EncodeToString([]byte(content))
	}
	return Wire
}

func Encode(Headers []Header, fromStaticTable, fromHeaderTable, toHuffman bool, table *Table, headerTableSize int) (Wire string) {
	if headerTableSize != -1 {
		intRep := PackIntRepresentation(uint32(headerTableSize), 5)
		(*intRep)[0] |= 0x20
		Wire += hex.EncodeToString(*intRep)
	}

	for _, header := range Headers {
		match, index := table.FindHeader(header)
		if fromStaticTable && match {
			var indexLen, mask byte
			var content string
			if fromHeaderTable {
				indexLen = 7
				mask = 0x80
				content = ""
			} else {
				indexLen = 4
				mask = 0x00
				content = PackContent(header.Value, toHuffman)
			}
			intRep := PackIntRepresentation(uint32(index), indexLen)
			(*intRep)[0] |= mask
			Wire += hex.EncodeToString(*intRep) + content
		} else if fromStaticTable && !match && index > 0 {
			var indexLen, mask byte
			if fromHeaderTable {
				indexLen = 6
				mask = 0x40
				table.AddHeader(header)
			} else {
				indexLen = 4
				mask = 0x00
			}
			intRep := PackIntRepresentation(uint32(index), indexLen)
			(*intRep)[0] |= mask
			Wire += hex.EncodeToString(*intRep) + PackContent(header.Value, toHuffman)
		} else {
			var prefix string
			if fromHeaderTable {
				prefix = "40"
				table.AddHeader(header)
			} else {
				prefix = "00"
			}
			content := PackContent(header.Name, toHuffman) + PackContent(header.Value, toHuffman)
			Wire += prefix + content
		}
	}

	return
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

func ParseFromByte(buf []byte) (content string, cursor uint32) {
	length, cursor := ParseIntRepresentation(buf, 7)

	if buf[0]&0x80 > 0 {
		content = huffman.Root.Decode(buf[cursor:], length)
	} else {
		content = string(buf[cursor : cursor+length])
	}

	cursor += length
	return
}

func ParseHeader(index uint32, buf []byte, isIndexed bool, table *Table) (name, value string, cursor uint32) {
	if c := uint32(0); !isIndexed {
		if index == 0 {
			name, c = ParseFromByte(buf[cursor:])
			cursor += c
		}
		value, c = ParseFromByte(buf[cursor:])
		cursor += c
	}

	if index > 0 {
		header := table.GetHeader(index)

		name = header.Name
		if len(value) == 0 {
			value = header.Value
		}
	}
	return
}

func Decode(wire string, table *Table) (Headers []Header) {
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
			size, c := ParseIntRepresentation((*buf)[cursor:], 5)
			table.SetHeaderTableSize(size)
			cursor += c
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

		name, value, c := ParseHeader(index, (*buf)[cursor:], isIndexed, table)
		cursor += c

		header := Header{name, value}
		if isIncremental {
			table.AddHeader(header)
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
	//fmt.Println(Decode("00073a6d6574686f640347455400073a736368656d650468747470000a3a617574686f726974790f7777772e7961686f6f2e636f2e6a7000053a70617468012f"))
	fmt.Println(huffman.HUFFMAN_TABLE)
}
