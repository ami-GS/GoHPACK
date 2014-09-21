package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
			fmt.Printf("%v\n", jsontype)
		}
	}
}
