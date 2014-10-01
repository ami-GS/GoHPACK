package hpack

type Header struct {
	Name, Value string
}

func (h Header) size() uint32 {
	return uint32(len(h.Name + h.Value))
}

func FindHeader(name, value string) (bool, int) {
	h := Header{name, value}
	for i, header := range *STATIC_TABLE {
		if header == h {
			return true, i
		} else if header.Name == name {
			return false, i
		}
	}

	ring := head
	for i := 0; i < int(currentEntryNum); i++ {
		if ring.header == h {
			return true, i + int(STATIC_TABLE_NUM)
		} else if ring.header.Name == name {
			return false, i + int(STATIC_TABLE_NUM)
		}
		ring = ring.Next
	}
	return false, -1

}

func GetHeader(index uint32) Header {
	if 0 < index && index < uint32(STATIC_TABLE_NUM) {
		return (*STATIC_TABLE)[index]
	} else if uint32(STATIC_TABLE_NUM) <= index && index <= uint32(STATIC_TABLE_NUM+byte(currentEntryNum)) {
		return getFromHeaderTable(index) //from Header Table
	} else {
		panic("error")
	}
}

type RingTable struct {
	header    Header
	Next, Pre *RingTable
}

var ringTable RingTable
var head *RingTable = &ringTable //*RingTalbe = &RingTable{Header{"", ""}, ringTable, nil}
var tail *RingTable = &ringTable //*RingTable = &RingTable{Header{"", ""}, nil, ringTable}

var HeaderTable *[]Header
var currentEntryNum uint16 = 0
var currentEntrySize uint32 = 0
var headerTableSize uint32 = 4096

func getFromHeaderTable(index uint32) Header {
	//return (*HeaderTable)[index-uint32(STATIC_TABLE_NUM)]
	index -= uint32(STATIC_TABLE_NUM)
	ring := head.Next
	for i := uint32(0); i < index; i++ {
		ring = ring.Next
	}
	return ring.header
}

var nilElem *RingTable

func delLast() {
	currentEntrySize -= tail.Pre.header.size()
	tail.Pre = tail.Pre.Pre
	currentEntryNum--
}

func insertFirst(header Header) {
	//here should be refactored
	elem := RingTable{header, head.Next, nil}
	if currentEntryNum >= 1 {
		head.Next.Pre = &elem
	}
	head.Next = &elem

	if currentEntryNum == 1 {
		tail.Pre = &elem
	}
	currentEntrySize += header.size()
	currentEntryNum++
}

func AddHeader(header Header) {
	for currentEntrySize+header.size() > headerTableSize {
		delLast()
	}
	insertFirst(header)
}

func SetMaxHeaderTableSize(size uint32) {
	headerTableSize = size
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
