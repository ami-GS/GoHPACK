package huffman

import (
	"reflect"
	"testing"
)

//TODO more precise test cases should be applied

func TestEncode(t *testing.T) {
	Root := Node{nil, nil, -1}
	Root.CreateTree()
	actualS, actualL := Root.Encode("http://www.amazon.com/")
	expectedS := []byte{157, 41, 174, 227, 12, 120, 241, 225, 113, 210, 63, 103, 169, 114, 30, 150, 63}
	expectedL := uint32(17)
	if !reflect.DeepEqual(actualS, expectedS) {
		t.Errorf("got %v\nwant %v", actualS, expectedS)
	}
	if actualL != expectedL {
		t.Errorf("got %v\nwant%v", actualL, expectedL)
	}
}

func TestDecode(t *testing.T) {
	Root := Node{nil, nil, -1}
	Root.CreateTree()
	actual := Root.Decode([]byte{157, 41, 174, 227, 12, 120, 241, 225, 113, 210, 63, 103, 169, 114, 30, 150, 63}, 17)
	expected := "http://www.amazon.com/"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
