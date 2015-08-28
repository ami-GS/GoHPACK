package GoHPACK

import (
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
		}
		return []byte{0x00}
	}

	var Wire []byte
	if toHuffman {

		encoded, length := huffman.Root.Encode(content)
		intRep := PackIntRepresentation(uint32(length), 7)
		intRep[0] |= 0x80

		//Wire += hex.EncodeToString(*intRep) + strings.Trim(hex.EncodeToString(b), "00") // + encoded
		return append(append(Wire, intRep...), encoded...)
	}

	intRep := PackIntRepresentation(uint32(len(content)), 7)
	return append(append(Wire, intRep...), []byte(content)...)
}

func Encode(Headers []Header, fromStaticTable, fromDynamicTable, toHuffman bool, table *Table, dynamicTableSize int) (Wire []byte) {
	if dynamicTableSize != -1 {
		Wire = PackIntRepresentation(uint32(dynamicTableSize), 5)
		Wire[0] |= 0x20
	}

	for _, header := range Headers {
		match, index := table.FindHeader(header)
		if fromStaticTable && match {
			if fromDynamicTable {
				intRep := PackIntRepresentation(uint32(index), 7)
				intRep[0] |= 0x80
				Wire = append(Wire, intRep...)
			} else {
				intRep := PackIntRepresentation(uint32(index), 4)
				Wire = append(append(Wire, intRep...),
					PackContent(header.Value, toHuffman)...)
			}
		} else if fromStaticTable && !match && index > 0 {
			var intRep []byte
			if fromDynamicTable {
				intRep = PackIntRepresentation(uint32(index), 6)
				intRep[0] |= 0x40
				table.AddHeader(header)
			} else {
				intRep = PackIntRepresentation(uint32(index), 4)
			}
			Wire = append(append(Wire, intRep...),
				PackContent(header.Value, toHuffman)...)
		} else {
			var prefix byte = 0x00
			if fromDynamicTable {
				prefix = 0x40
				table.AddHeader(header)
			}
			content := append(PackContent(header.Name, toHuffman), PackContent(header.Value, toHuffman)...)
			Wire = append(append(Wire, prefix), content...)
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
		if buf[cursor]&0xe0 == 0x20 {
			// 7.3 Header Table Size Update
			size, c := ParseIntRepresentation(buf[cursor:], 5)
			table.SetDynamicTableSize(size)
			cursor += c
		}

		var nLen byte
		if (buf[cursor] & 0x80) > 0 {
			// 7.1 Indexed Header Field
			if (buf[cursor] & 0x7f) == 0 {
				panic('a')
			}
			nLen = 7
			isIndexed = true
		} else {
			if buf[cursor]&0xc0 == 0x40 {
				// 7.2.1 Literal Header Field with Incremental Indexing
				nLen = 6
			} else {
				// when buf[cursor]&0xf0 == 0xf0
				// 7.2.2 Literal Header Field without Indexing
				// else
				// 7.2.3 Literal Header Field never Indexed
				nLen = 4
			}
		}
		index, c1 := ParseIntRepresentation(buf[cursor:], nLen)
		name, value, c2 := ParseHeader(index, buf[cursor+c1:], isIndexed, table)
		cursor += c1 + c2
		header := Header{name, value}
		if nLen == 6 {
			table.AddHeader(header)
		}
		Headers = append(Headers, header)
	}

	return
}
