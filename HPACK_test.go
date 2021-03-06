package GoHPACK

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

var N uint32 = 10000000

type jsonobject struct {
	Cases       []Case
	Draft       int
	Description string
}

type Case struct {
	Seqno             int
	Header_table_size uint32
	Wire              string
	Headers           []map[string]string
}

var TESTCASES = []string{
	"hpack-test-case/haskell-http2-naive/",
	"hpack-test-case/haskell-http2-naive-huffman/",
	"hpack-test-case/haskell-http2-static/",
	"hpack-test-case/haskell-http2-static-huffman/",
	"hpack-test-case/haskell-http2-linear/",
	"hpack-test-case/haskell-http2-linear-huffman/",
	"hpack-test-case/go-hpack/",
	"hpack-test-case/nghttp2/",
	"hpack-test-case/nghttp2-16384-4096/",
	"hpack-test-case/nghttp2-change-table-size/",
	"hpack-test-case/node-http2-hpack/",
}

/*
func TestPackIntRepresentation(t *testing.T) {
	// TODO: need test cases
	var i uint32 = 0
	var n byte = 0
	for ; i < N; i++ {
		n = 1
		for ; n < 8; n++ {
			actual := PackIntRepresentation(i, n)
		}
	}
}
*/
func ConvertHeader(headers []map[string]string) (dist []Header) {
	for _, dict := range headers {
		for k, v := range dict {
			dist = append(dist, Header{k, v})
		}
	}
	return
}

func TestDecode(t *testing.T) {
	for _, testCase := range TESTCASES {
		files, err := ioutil.ReadDir(testCase)
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			table := NewTable()
			data, err := ioutil.ReadFile(testCase + file.Name())
			if err != nil {
				panic(err)
			}
			var jsontype jsonobject
			json.Unmarshal(data, &jsontype)

			for _, seq := range jsontype.Cases {
				buf, err := hex.DecodeString(seq.Wire)
				if err != nil {
					panic(err)
				}
				actual := Decode(buf, &table)
				expected := ConvertHeader(seq.Headers)
				if !reflect.DeepEqual(actual, expected) {
					t.Errorf("get %v\nwant %v", actual, expected)
					t.Errorf("False in %s at seq %d", testCase+file.Name(), seq.Seqno)
				}
			}
		}
	}
}

func EncType(testCase string) (fStatic, fHeader, isHuffman bool) {
	fHeader = strings.Contains(testCase, "linear")
	if fHeader {
		fStatic = true
	} else {
		fStatic = strings.Contains(testCase, "static")
	}
	if strings.Contains(testCase, "haskell") {
		isHuffman = strings.Contains(testCase, "huffman")
	} else {
		isHuffman = true
	}
	return
}

func TestEncode(t *testing.T) {
	for _, testCase := range TESTCASES {
		fStatic, fHeader, isHuffman := EncType(testCase)
		files, err := ioutil.ReadDir(testCase)
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			encTable := NewTable()
			decTable := NewTable()
			data, err := ioutil.ReadFile(testCase + file.Name())
			if err != nil {
				panic(err)
			}
			var jsontype jsonobject
			json.Unmarshal(data, &jsontype)

			for _, seq := range jsontype.Cases {
				var d_table_size int = int(seq.Header_table_size)
				if d_table_size == 0 {
					d_table_size = -1
				}
				expected := ConvertHeader(seq.Headers)
				buf := Encode(expected, fStatic, fHeader, isHuffman, &encTable, d_table_size)
				actual := Decode(buf, &decTable)
				if !reflect.DeepEqual(actual, expected) {
					t.Errorf("get %v\nwant %v", actual, expected)
					t.Errorf("False in %s at seq %d", testCase+file.Name(), seq.Seqno)
				}
			}
		}
	}
}
