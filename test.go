package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	hpack "github.com/ami-GS/GoHPACK"
	"github.com/ami-GS/GoHPACK/huffman"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

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

var TESTCASE = []string{
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

func EncType(testCase string) (fStatic, fHeader, isHuffman bool) {
	fHeader = strings.Contains(testCase, "linear")
	if fHeader {
		fStatic = true
	} else {
		fStatic = strings.Contains(testCase, "static")
	}
	isHuffman = strings.Contains(testCase, "huffman")
	return
}

func convertHeader(headers []map[string]string) (dist []hpack.Header) {
	for _, dict := range headers {
		for k, v := range dict {
			dist = append(dist, hpack.Header{k, v})
		}
	}
	return
}

func compHeaders(decoded, correct []hpack.Header, storyPass *bool) {
	if !reflect.DeepEqual(correct, decoded) {
		*storyPass = false
		if len(os.Args) == 3 && os.Args[2] == "-v" {
			fmt.Println(decoded)
			fmt.Println(correct)
		}
		//os.Exit(-1)
	}
}

func compWire(encoded, correct string, storyPass *bool) {
	if encoded != correct {
		*storyPass = false
		if len(os.Args) == 3 && os.Args[2] == "-v" {
			fmt.Println(encoded)
			fmt.Println(correct)
		}
		//os.Exit(-1)
	}

}

func main() {
	fmt.Println(hex.EncodeToString(*hpack.PackIntRepresentation(3000000, 5)))
	nums, _ := hex.DecodeString(string("1fa18db701"))
	fmt.Println(hpack.ParseIntRepresentation(nums, 5))
	huffman.Root.CreateTree()
	for _, testCase := range TESTCASE {
		fStatic, fHeader, isHuffman := EncType(testCase)
		files, err := ioutil.ReadDir(testCase)
		if err != nil {
			panic(err)
		}
		for _, f := range files {
			encTable := hpack.InitTable()
			decTable := hpack.InitTable()
			data, err := ioutil.ReadFile(testCase + f.Name())
			if err != nil {
				panic(err)
			}
			var jsontype jsonobject
			json.Unmarshal(data, &jsontype)
			storyPass := true

			if len(os.Args) >= 2 && os.Args[1] == "-d" {
				for _, seq := range jsontype.Cases {
					Headers := hpack.Decode(seq.Wire, &decTable)
					correctHeaders := convertHeader(seq.Headers)
					compHeaders(Headers, correctHeaders, &storyPass)
					if !storyPass {
						fmt.Println("False in", testCase+f.Name(), "at seq", seq.Seqno)
						break
					}

				}
			} else if len(os.Args) >= 2 && (os.Args[1] == "-e" || os.Args[1] == "-a") {
				for _, seq := range jsontype.Cases {
					if seq.Header_table_size != 0 {
						encTable.SetHeaderTableSize(seq.Header_table_size)
					}
					Headers := convertHeader(seq.Headers)
					Wire := hpack.Encode(Headers, fStatic, fHeader, isHuffman, &encTable, -1)
					if os.Args[1] == "-a" {
						distHeaders := hpack.Decode(Wire, &decTable)
						compHeaders(distHeaders, Headers, &storyPass)
						if !storyPass {
							fmt.Println("False in", testCase+f.Name(), "at seq", seq.Seqno)
							break
						}
					} else {
						compWire(Wire, seq.Wire, &storyPass)
						if !storyPass {
							fmt.Println("False in", testCase+f.Name(), "at seq", seq.Seqno)
							break
						}
					}
				}
			} else {
				fmt.Println("argument should be '-e', '-d', '-a' or none")
				os.Exit(-1)
			}
			if storyPass {
				fmt.Println("Pass in " + testCase + f.Name())
			}
		}
	}
}
