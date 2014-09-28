package main

import (
	"encoding/json"
	"fmt"
	"hpack"
	"io/ioutil"
	"os"
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
			storyPass := true
			for _, seq := range jsontype.Cases {
				//fmt.Printf("%d\n", len(seq.Wire))
				Headers := hpack.Decode(seq.Wire)
				testHeaders := []hpack.Header{}
				for _, dict := range seq.Headers {
					for k, v := range dict {
						testHeaders = append(testHeaders, hpack.Header{k, v})
					}
				}

				if !reflect.DeepEqual(testHeaders, Headers) {
					storyPass = false
					fmt.Println("False in", f.Name(), "at seq", seq.Seqno)
					fmt.Println(Headers)
					fmt.Println(testHeaders)
					os.Exit(-1)
					break
				}
			}
			if storyPass {
				fmt.Println("Pass in", f.Name())
			}
		}
	}
}
