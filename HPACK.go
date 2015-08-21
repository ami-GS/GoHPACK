//package main //for test
package GoHPACK

import (
	"fmt"
	"github.com/ami-GS/GoHPACK/huffman"
)

func PackIntRepresentation(I uint32, N byte) (buf []byte) {
	if I < uint32(1<<N)-1 {
		return []byte{byte(I)}
	}

	I -= uint32(1<<N) - 1
	var i int = 1
	tmpI := I
	for ; tmpI >= 128; i++ {
		tmpI = tmpI >> 7
	} // check length

	buf = make([]byte, i+1)
	buf[0] = byte(1<<N) - 1
	i = 1
	for ; I >= 0x80; i++ {
		buf[i] = (byte(I) & 0x7f) | 0x80
		I = I >> 7
	}
	buf[i] = byte(I)

	return buf

}

func PackContent(content string, toHuffman bool) []byte {
	if len(content) == 0 {
		if toHuffman {
			return []byte{0x80}
		} else {
			return []byte{0x00}
		}
	}

	var Wire []byte
	if toHuffman {

		encoded, length := huffman.Root.Encode(content)
		intRep := PackIntRepresentation(uint32(length), 7)
		intRep[0] |= 0x80

		//Wire += hex.EncodeToString(*intRep) + strings.Trim(hex.EncodeToString(b), "00") // + encoded
		Wire = append(append(Wire, intRep...), encoded...)
	} else {
		intRep := PackIntRepresentation(uint32(len(content)), 7)
		Wire = append(append(Wire, intRep...), []byte(content)...)
	}
	return Wire
}

func Encode(Headers []Header, fromStaticTable, fromDynamicTable, toHuffman bool, table *Table, dynamicTableSize int) (Wire []byte) {
	if dynamicTableSize != -1 {
		intRep := PackIntRepresentation(uint32(dynamicTableSize), 5)
		intRep[0] |= 0x20
		Wire = intRep
	}

	for _, header := range Headers {
		match, index := table.FindHeader(header)
		if fromStaticTable && match {
			var indexLen, mask byte
			var content []byte
			if fromDynamicTable {
				indexLen = 7
				mask = 0x80
				content = []byte{}
			} else {
				indexLen = 4
				mask = 0x00
				content = PackContent(header.Value, toHuffman)
			}
			intRep := PackIntRepresentation(uint32(index), indexLen)
			intRep[0] |= mask
			Wire = append(append(Wire, intRep...), content...)
		} else if fromStaticTable && !match && index > 0 {
			var indexLen, mask byte
			if fromDynamicTable {
				indexLen = 6
				mask = 0x40
				table.AddHeader(header)
			} else {
				indexLen = 4
				mask = 0x00
			}
			intRep := PackIntRepresentation(uint32(index), indexLen)
			intRep[0] |= mask
			Wire = append(append(Wire, intRep...), PackContent(header.Value, toHuffman)...)
		} else {
			var prefix []byte
			if fromDynamicTable {
				prefix = []byte{0x40}
				table.AddHeader(header)
			} else {
				prefix = []byte{0x00}
			}
			content := append(PackContent(header.Name, toHuffman), PackContent(header.Value, toHuffman)...)
			Wire = append(append(Wire, prefix...), content...)
		}
	}

	return
}

func ParseIntRepresentation(buf []byte, N byte) (I, cursor uint32) {
	I = uint32(buf[0] & ((1 << N) - 1)) // byte could be used as byte
	cursor = 1
	if I < ((1 << N) - 1) {
		return I, cursor
	}

	var M byte = 0
	for (buf[cursor] & 0x80) > 0 {
		I += uint32(buf[cursor]&0x7f) * (1 << M)
		M += 7
		cursor += 1
	}
	I += uint32(buf[cursor]&0x7f) * (1 << M)
	return I, cursor + 1

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

func Decode(buf []byte, table *Table) (Headers []Header) {
	var cursor uint32 = 0
	for cursor < uint32(len(buf)) {
		isIndexed := false
		isIncremental := false
		var index, c uint32
		if buf[cursor]&0xe0 == 0x20 {
			// 7.3 Header Table Size Update
			size, c := ParseIntRepresentation(buf[cursor:], 5)
			table.SetDynamicTableSize(size)
			cursor += c
		}

		if (buf[cursor] & 0x80) > 0 {
			// 7.1 Indexed Header Field
			if (buf[cursor] & 0x7f) == 0 {
				panic('a')
			}
			index, c = ParseIntRepresentation(buf[cursor:], 7)
			isIndexed = true
		} else {
			if buf[cursor]&0xc0 == 0x40 {
				// 7.2.1 Literal Header Field with Incremental Indexing
				index, c = ParseIntRepresentation(buf[cursor:], 6)
				isIncremental = true
			} else if buf[cursor]&0xf0 == 0xf0 {
				// 7.2.3 Literal Header Field never Indexed
				index, c = ParseIntRepresentation(buf[cursor:], 4)
			} else {
				// 7.2.2 Literal Header Field without Indexing
				index, c = ParseIntRepresentation(buf[cursor:], 4)
			}
		}
		cursor += c

		name, value, c := ParseHeader(index, buf[cursor:], isIndexed, table)
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
