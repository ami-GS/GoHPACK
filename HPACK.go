package GoHPACK

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

func Encode(Headers []Header, fromStaticTable, fromDynamicTable, toHuffman bool, table *Table, dynamicTableSize int) (Wire []byte) {
	if dynamicTableSize != -1 && int(table.dynamicTableSize) != dynamicTableSize {
		table.SetDynamicTableSize(uint32(dynamicTableSize))
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
					table.PackContent(header.Value, toHuffman)...)
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
				table.PackContent(header.Value, toHuffman)...)
		} else {
			var prefix byte = 0x00
			if fromDynamicTable {
				prefix = 0x40
				table.AddHeader(header)
			}
			content := append(table.PackContent(header.Name, toHuffman),
				table.PackContent(header.Value, toHuffman)...)
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
		name, value, c2 := table.ParseHeader(index, buf[cursor+c1:], isIndexed)
		cursor += c1 + c2
		header := Header{name, value}
		if nLen == 6 {
			table.AddHeader(header)
		}
		Headers = append(Headers, header)
	}

	return
}
