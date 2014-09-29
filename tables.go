package hpack

type Header struct {
	Name, Value string
}

func FindHeader(name, value string) (bool, int) {
	h := Header{name, value}
	for i, header := range *STATIC_TABLE {
		if header == h {
			return true, i
		} else if header.Name == name {
			return false, i
		} else {
			return false, -1
		}
	}
	return false, -1
}

var currentEntryNum byte = 0

func GetHeader(index uint32) Header {
	if 0 < index && index < uint32(STATIC_TABLE_NUM) {
		return (*STATIC_TABLE)[index]
	} else if uint32(STATIC_TABLE_NUM) <= index && index <= uint32(STATIC_TABLE_NUM+currentEntryNum) {
		return Header{"", ""} //from Header Table
	} else {
		panic("error")
	}
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

//var NAME_TABLE = [header[0} for header in STATIC_TABLE]
