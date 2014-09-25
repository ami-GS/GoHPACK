package main

import (
	"encoding/json"
	"fmt"
	"hpack"
	"io/ioutil"
	"reflect"
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

type Header struct {
	Name, Value string
}

var TESTCASE = []string{
	"hpack-test-case/haskell-http2-naive/",
	//"hpack-test-case/haskell-http2-naive-huffman",
}

func main() {
	for _, testCase := range TESTCASE {
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
			for _, seq := range jsontype.Cases {
				Headers := hpack.Decode(seq.Wire)
				if reflect.DeepEqual(seq.Headers, Headers) {
					fmt.Println("Pass in", f.Name())
				} else {
					fmt.Println("False in", f.Name(), "at seq", seq.Seqno)
					break
				}
			}
		}
	}
}
