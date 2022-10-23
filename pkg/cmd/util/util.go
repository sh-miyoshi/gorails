package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunCommand(name string, args ...string) {
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		fmt.Printf("Failed to run %s %v: %+v", name, args, err)
		os.Exit(1)
	}
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func AppendLine(fname string, data string) {
	const marker = "GORAILS MARKER"

	fp, err := os.Open(fname)
	if err != nil {
		fmt.Printf("Failed to open file %s: %+v", fname, err)
		os.Exit(1)
	}

	results := []string{}
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		results = append(results, line)

		if strings.Contains(line, marker) {
			results = append(results, data)
		}
	}
	fp.Close()

	// Write to file
	fp, err = os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Printf("Failed to open file %s for writing: %+v", fname, err)
		os.Exit(1)
	}

	for _, line := range results {
		fp.WriteString(line + "\n")
	}
	fp.Close()
}

func CamelToSnake(s string) string {
	if s == "" {
		return s
	}

	delimiter := "_"
	sLen := len(s)
	var snake string
	for i, current := range s {
		if i > 0 && i+1 < sLen {
			if current >= 'A' && current <= 'Z' {
				next := s[i+1]
				prev := s[i-1]
				if (next >= 'a' && next <= 'z') || (prev >= 'a' && prev <= 'z') {
					snake += delimiter
				}
			}
		}
		snake += string(current)
	}

	snake = strings.ToLower(snake)
	return snake
}
