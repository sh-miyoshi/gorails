package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

func RunCommand(name string, args ...string) {
	c := exec.Command(name, args...)
	out, err := c.Output()
	if err != nil {
		fmt.Printf("Failed to run %s %v: %+v", name, args, err)
		os.Exit(1)
	}
	if len(out) > 0 {
		fmt.Println(string(out))
	}
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func CopyTemplateFile(src, dst string, data any) {
	tpl, err := template.ParseFiles(src)
	if err != nil {
		fmt.Printf("Failed to parse template %s: %+v", src, err)
		os.Exit(1)
	}
	fp, err := os.Create(dst)
	if err != nil {
		fmt.Printf("Failed to create new file %s: %+v", dst, err)
		os.Exit(1)
	}
	defer fp.Close()

	tpl.Execute(fp, data)
}

func CopyFile(src, dst string) {
	srcFp, err := os.Open(src)
	if err != nil {
		fmt.Printf("Failed to open source file %s: %+v", src, err)
		os.Exit(1)
	}
	defer srcFp.Close()

	dstFp, err := os.Create(dst)
	if err != nil {
		fmt.Printf("Failed to create distination file %s: %+v", dst, err)
		os.Exit(1)
	}
	defer dstFp.Close()

	io.Copy(dstFp, srcFp)
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
