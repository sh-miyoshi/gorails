package util

import (
	"os"
	"testing"
)

func TestAppendLine(t *testing.T) {
	// Case 1 a file with marker
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

	expect := `AAA
BBB
// GORAILS MARKER Don't edit this line
DDD
CCC
`
	buf, _ := os.ReadFile(f1.Name())
	if string(buf) != expect {
		t.Errorf("Failed to append data, expect %s, but got %s", expect, string(buf))
	}

	// Case 2 a file with no marker
	f2, _ := os.CreateTemp("", "")
	defer os.Remove(f2.Name())

	data = `AAA
BBB
CCC
`
	f2.WriteString(data)
	f2.Close()
	AppendLine(f2.Name(), "DDD")

	expect = `AAA
BBB
CCC
`
	buf, _ = os.ReadFile(f2.Name())
	if string(buf) != expect {
		t.Errorf("Failed to append data, expect %s, but got %s", expect, string(buf))
	}
}
