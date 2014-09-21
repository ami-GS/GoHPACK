package main //for test

import (
	"encoding/hex"
	"fmt"
)

type Header struct {
	Name, Value string
}

func decode(wire string) []Header {
	var Headers []Header
	var buf *[]byte
	nums, err := hex.DecodeString(string(wire))
	if err != nil {
		panic(err)
	}
	buf = &nums

	/* Usage of (*buf)[]
	for i := 0; i < 5; i++ {
		fmt.Println((*buf)[i])
	}
	*/

	for i, v := range *buf {
		fmt.Println(i, v)
	}

	d := hex.EncodeToString(nums)
	fmt.Println(d)
	return Headers
}

func main() {
	decode("ff80000111")
	fmt.Println()
}
