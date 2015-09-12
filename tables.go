package GoHPACK

import (
	"github.com/ami-GS/GoHPACK/huffman"
)

type Header struct {
	Name, Value string
}

func (h Header) size() uint32 {
	return uint32(len(h.Name + h.Value))
}

type Table struct {
	head, tail       *RingTable
	currentEntrySize uint32
	currentEntryNum  uint32
	dynamicTableSize uint32
	Huffman          *huffman.Node
}

type RingTable struct {
	header    Header
	Next, Pre *RingTable
}

func NewTable() (t Table) {
	t.currentEntryNum = 0
	t.currentEntrySize = 0
	t.dynamicTableSize = 4096
	t.Huffman = &huffman.Node{nil, nil, -1} // temporally
	t.Huffman.CreateTree()
	return
}

func (t *Table) FindHeader(h Header) (match bool, index int) {
	//here should be optimized
	preName := ""
	for i, header := range STATIC_TABLE {
		if header == h {
			return true, i
		} else if header.Name == h.Name && index == 0 {
			preName = header.Name
			index = i
			match = false
		} else if index != 0 && preName != header.Name {
			return match, index
		}
	}

	ring := t.head
	for i := 0; i < int(t.currentEntryNum); i++ {
		if ring.header == h {
			return true, i + int(STATIC_TABLE_NUM)
		} else if ring.header.Name == h.Name && index == 0 {
			match = false
			index = i + int(STATIC_TABLE_NUM)
		}
		ring = ring.Next
	}
	if index > 0 {
		return match, index
	}
	return false, -1 // not found on any table
}

func (t *Table) GetHeader(index uint32) Header {
	if 0 < index && index < uint32(STATIC_TABLE_NUM) {
		return STATIC_TABLE[index]
	} else if uint32(STATIC_TABLE_NUM) <= index && index <= uint32(STATIC_TABLE_NUM+byte(t.currentEntryNum)) {
		return t.getFromDynamicTable(index) //from Header Table
	}

	panic("error")
}

func (t *Table) getFromDynamicTable(index uint32) Header {
	index -= uint32(STATIC_TABLE_NUM)
	ring := t.head
	for i := uint32(0); i < index; i++ {
		ring = ring.Next
	}
	return ring.header
}

var nilElem *RingTable

func (t *Table) delLast() {
	t.currentEntrySize -= t.tail.header.size()
	deleated := t.tail
	t.tail = deleated.Pre
	deleated = nilElem
	t.currentEntryNum--
}

func (t *Table) insertFirst(header Header) {
	//here should be refactored
	elem := RingTable{header, nil, nil}

	if t.currentEntryNum >= 1 {
		elem.Next = t.head
		t.head.Pre = &elem
	}
	t.head = &elem

	if t.currentEntryNum == 0 {
		t.tail = &elem
	}
	t.currentEntryNum++
	t.currentEntrySize += header.size()
}

func (t *Table) AddHeader(header Header) {
	for t.currentEntrySize+header.size() > t.dynamicTableSize {
		t.delLast()

func (t *Table) PackContent(content string, toHuffman bool) []byte {
	if toHuffman {
		encoded, length := t.Huffman.Encode(content)
		intRep := PackIntRepresentation(uint32(length), 7)
		intRep[0] |= 0x80
		return append(intRep, encoded...)
	}
	intRep := PackIntRepresentation(uint32(len(content)), 7)
	return append(intRep, []byte(content)...)
}

func (t *Table) SetDynamicTableSize(size uint32) {
	t.dynamicTableSize = size
func (t *Table) ParseFromByte(buf []byte) (content string, cursor uint32) {
	length, cursor := ParseIntRepresentation(buf, 7)

	if buf[0]&0x80 > 0 {
		content = t.Huffman.Decode(buf[cursor:], length)
	} else {
		content = string(buf[cursor : cursor+length])
	}
	cursor += length
	return
}

func (table *Table) ParseHeader(index uint32, buf []byte, isIndexed bool) (name, value string, cursor uint32) {
	if c := uint32(0); !isIndexed {
		if index == 0 {
			name, c = table.ParseFromByte(buf[cursor:])
			cursor += c
		}
		value, c = table.ParseFromByte(buf[cursor:])
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

var STATIC_TABLE = &[...]Header{
	{"", ""},
	{":authority", ""},
	{":method", "GET"},
	{":method", "POST"},
	{":path", "/"},
	{":path", "/index.html"},
	{":scheme", "http"},
	{":scheme", "https"},
	{":status", "200"},
	{":status", "204"},
	{":status", "206"},
	{":status", "304"},
	{":status", "400"},
	{":status", "404"},
	{":status", "500"},
	{"accept-charset", ""},
	{"accept-encoding", "gzip, deflate"},
	{"accept-language", ""},
	{"accept-ranges", ""},
	{"accept", ""},
	{"access-control-allow-origin", ""},
	{"age", ""},
	{"allow", ""},
	{"authorization", ""},
	{"cache-control", ""},
	{"content-disposition", ""},
	{"content-encoding", ""},
	{"content-language", ""},
	{"content-length", ""},
	{"content-location", ""},
	{"content-range", ""},
	{"content-type", ""},
	{"cookie", ""},
	{"date", ""},
	{"etag", ""},
	{"expect", ""},
	{"expires", ""},
	{"from", ""},
	{"host", ""},
	{"if-match", ""},
	{"if-modified-since", ""},
	{"if-none-match", ""},
	{"if-range", ""},
	{"if-unmodified-since", ""},
	{"last-modified", ""},
	{"link", ""},
	{"location", ""},
	{"max-forwards", ""},
	{"proxy-authenticate", ""},
	{"proxy-authorization", ""},
	{"range", ""},
	{"referer", ""},
	{"refresh", ""},
	{"retry-after", ""},
	{"server", ""},
	{"set-cookie", ""},
	{"strict-transport-security", ""},
	{"transfer-encoding", ""},
	{"user-agent", ""},
	{"vary", ""},
	{"via", ""},
	{"www-authenticate", ""},
}
var STATIC_TABLE_NUM byte = byte(len(STATIC_TABLE))
