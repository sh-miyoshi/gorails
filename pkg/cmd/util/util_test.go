package util

import (
	"os"
	"testing"
)

func TestAppendLine(t *testing.T) {
	// Case 1 with marker
	f1, _ := os.CreateTemp("", "")
	defer os.Remove(f1.Name())

	data := `AAA
BBB
// GORAILS MARKER Don't edit this line
CCC
`
	f1.WriteString(data)
	f1.Close()
	AppendLine(f1.Name(), "DDD")

	t.Error(f1.Name())

	// expect := `AAA
	// BBB
	// // GORAILS MARKER Don't edit this line
	// DDD
	// CCC
	// `
	// buf := []byte{}
	// fp, _ := os.Open(f1.Name())
	// fp.Read(buf)
	// fp.Close()

	// if string(buf) != expect {
	// 	t.Errorf("Failed to append data, expect %s, but got %s", expect, string(buf))
	// }

	// create tmp file with data (no marker)
	// run, read file and check data
}
