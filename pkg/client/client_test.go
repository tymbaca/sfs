package sfs

import (
	"bytes"
	"io"
	"testing"
)

func Test_split(t *testing.T) {
	// t.Run("1234512345123", func(t *testing.T) {
	// 	r := strings.NewReader("1234512345123")
	//
	// 	rs, err := split(r, 5)
	// 	if err != nil {
	// 		t.FailNow()
	// 	}
	//
	// 	if len(rs) != 3 {
	// 		t.FailNow()
	// 	}
	//
	// 	assertReaderString(t, rs[0], "12345")
	// 	assertReaderString(t, rs[1], "12345")
	// 	assertReaderString(t, rs[2], "123")
	// })
	// t.Run("empty", func(t *testing.T) {
	// 	r := strings.NewReader("")
	//
	// 	rs, err := split(r, 5)
	// 	if err != nil {
	// 		t.FailNow()
	// 	}
	//
	// 	if len(rs) != 0 {
	// 		t.FailNow()
	// 	}
	// })
}

func assertReaderString(t *testing.T, r io.Reader, data string) {
	assertReader(t, r, []byte(data))
}
func assertReader(t *testing.T, r io.Reader, data []byte) {
	chunk, err := io.ReadAll(r)
	if err != nil {
		t.FailNow()
	}

	if !bytes.Equal(chunk, data) {
		t.FailNow()
	}
}
