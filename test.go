package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hpack"
	"huffman"
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
	Seqno   int
	Wire    string
	Headers []map[string]string
}

var TESTCASE = []string{
	"hpack-test-case/haskell-http2-naive/",
	"hpack-test-case/haskell-http2-naive-huffman/",
	"hpack-test-case/haskell-http2-static/",
	"hpack-test-case/haskell-http2-static-huffman/",
	"hpack-test-case/haskell-http2-linear/",
	"hpack-test-case/haskell-http2-linear-huffman/",
}

func main() {
	fmt.Println(hex.EncodeToString(*hpack.PackIntRepresentation(3000000, 5)))
	nums, _ := hex.DecodeString(string("1fa18db701"))
	fmt.Println(hpack.ParseIntRepresentation(nums, 5))
	huffman.Root.CreateTree()
	args := os.Args
	var fStatic, fHeader, isHuffman bool
	for _, testCase := range TESTCASE {
		fHeader = strings.Contains(testCase, "linear")
		if fHeader {
			fStatic = true
		} else {
			fStatic = strings.Contains(testCase, "static")
		}
		isHuffman = strings.Contains(testCase, "huffman")
		files, err := ioutil.ReadDir(testCase)
		if err != nil {
			panic(err)
		}
		for _, f := range files {
			data, err := ioutil.ReadFile(testCase + f.Name())
			if err != nil {
				panic(err)
			}
			var jsontype jsonobject
			json.Unmarshal(data, &jsontype)
			storyPass := true
			if len(args) == 1 || args[1] == "-d" {
				for _, seq := range jsontype.Cases {
					Headers := hpack.Decode(seq.Wire)
					testHeaders := []hpack.Header{}
					for _, dict := range seq.Headers {
						for k, v := range dict {
							testHeaders = append(testHeaders, hpack.Header{k, v})
						}
					}
					if !reflect.DeepEqual(testHeaders, Headers) {
						storyPass = false
						fmt.Println("False in", testCase+f.Name(), "at seq", seq.Seqno)
						fmt.Println(Headers)
						fmt.Println(testHeaders)
						os.Exit(-1)
						break
					}
				}
			} else if args[1] == "-e" {
				for _, seq := range jsontype.Cases {
					Headers := []hpack.Header{}
					for _, dict := range seq.Headers {
						for k, v := range dict {
							Headers = append(Headers, hpack.Header{k, v})
						}
					}
					Wire := hpack.Encode(Headers, fStatic, fHeader, isHuffman, -1)
					if Wire != seq.Wire {
						storyPass = false
						fmt.Println("False in", testCase+f.Name(), "at seq", seq.Seqno)
						fmt.Println(Wire)
						fmt.Println(seq.Wire)
						os.Exit(-1)
						break
					}
				}
			} else if args[1] == "-a" {
				for _, seq := range jsontype.Cases {
					Headers := []hpack.Header{}
					for _, dict := range seq.Headers {
						for k, v := range dict {
							Headers = append(Headers, hpack.Header{k, v})
						}
					}
					Wire := hpack.Encode(Headers, fStatic, fHeader, isHuffman, -1)
					distHeaders := hpack.Decode(Wire)
					if !reflect.DeepEqual(distHeaders, Headers) {
						storyPass = false
						fmt.Println("False in", testCase+f.Name(), "at seq", seq.Seqno)
						fmt.Println(distHeaders)
						fmt.Println(Headers)
						os.Exit(-1)
						break
					}
				}
			} else {
				fmt.Println("argument should be '-e' or '-d' or none")
				os.Exit(-1)
			}
			if storyPass {
				fmt.Println("Pass in " + testCase + f.Name())
			}
		}
	}
}
