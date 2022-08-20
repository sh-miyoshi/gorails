package util

import (
	"fmt"
	"io"
	"os"
	"os/exec"
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
