package util

import (
	"os"
	"testing"
)

func TestFileExists(t *testing.T) {
	// Case 1 check exists file
	file, _ := os.CreateTemp("", "")
	defer os.Remove(file.Name())
	if !FileExists(file.Name()) {
		t.Errorf("Case 1 failed. expect true, but got false")
	}

	// Case 2 check do not exist file or directory
	if FileExists("dummy") {
		t.Errorf("Case 2 failed. expect false, but got true")
	}

	// Case 3 check exists directory
	dir := os.TempDir()
	defer os.Remove(dir)
	if !FileExists(dir) {
		t.Errorf("Case 3 failed. expect true, but got false")
	}
}

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
