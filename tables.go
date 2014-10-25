package GoHPACK

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
	headerTableSize  uint32
}

type RingTable struct {
	header    Header
	Next, Pre *RingTable
}

func InitTable() (t Table) {
	t.currentEntryNum = 0
	t.currentEntrySize = 0
	t.headerTableSize = 4096
	return
}

func (t *Table) FindHeader(h Header) (match bool, index int) {
	//here should be optimized
	preName := ""
	for i, header := range *STATIC_TABLE {
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
	} else {
		return false, -1
	}
}

func (t *Table) GetHeader(index uint32) Header {
	if 0 < index && index < uint32(STATIC_TABLE_NUM) {
		return (*STATIC_TABLE)[index]
	} else if uint32(STATIC_TABLE_NUM) <= index && index <= uint32(STATIC_TABLE_NUM+byte(t.currentEntryNum)) {
		return t.getFromHeaderTable(index) //from Header Table
	} else {
		panic("error")
	}
}

func (t *Table) getFromHeaderTable(index uint32) Header {
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
	for t.currentEntrySize+header.size() > t.headerTableSize {
		t.delLast()
	}
	t.insertFirst(header)
}

func (t *Table) SetHeaderTableSize(size uint32) {
	t.headerTableSize = size
}

var STATIC_TABLE *[]Header = &[]Header{
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
var STATIC_TABLE_NUM byte = byte(len(*STATIC_TABLE))
