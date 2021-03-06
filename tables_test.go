package GoHPACK

import (
	"reflect"
	"testing"
)

func TestNewTable(t *testing.T) {
	actual := NewTable()
	expected := Table{
		head:             nil,
		tail:             nil,
		currentEntrySize: 0,
		currentEntryNum:  0,
		dynamicTableSize: 4096,
	}
	if actual.head != expected.head {
		t.Errorf("got %v\nwant %v", actual.head, expected.head)
	}
	if actual.tail != expected.tail {
		t.Errorf("got %v\nwant %v", actual.tail, expected.tail)
	}
	if actual.currentEntryNum != expected.currentEntryNum {
		t.Errorf("got %v\nwant %v", actual.currentEntryNum, expected.currentEntryNum)
	}
	if actual.currentEntrySize != expected.currentEntrySize {
		t.Errorf("got %v\nwant %v", actual.currentEntrySize, expected.currentEntrySize)
	}
	if actual.dynamicTableSize != expected.dynamicTableSize {
		t.Errorf("got %v\nwant %v", actual.dynamicTableSize, expected.dynamicTableSize)
	}
}

// TODO: DynamicTable should also be tested
func TestFindHeader(t *testing.T) {
	table := NewTable()
	for i, header := range STATIC_TABLE {
		actualB, actualI := table.FindHeader(header)
		expectedB := true
		expectedI := i
		if actualB != expectedB {
			t.Errorf("got %v\nwant %v", actualB, expectedB)
		}
		if actualI != expectedI {
			t.Errorf("got %v\nwant %v", actualI, expectedI)
		}
	}
}

func TestGetHeader(t *testing.T) {
	table := NewTable()
	for i := 1; byte(i) < STATIC_TABLE_NUM; i++ {
		actual := table.GetHeader(uint32(i))
		expected := STATIC_TABLE[i]
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("got %v\nwant %v", actual, expected)
		}
	}
}
