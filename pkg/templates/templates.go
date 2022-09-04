package templates

import (
	"fmt"
	"os"
)

const (
	GitIgnore int = iota
	ModelBase
)

func Exec(templateType int, dstFile string, data any) {
	fp, err := os.Create(dstFile)
	if err != nil {
		fmt.Printf("Failed to create new file %s: %+v", dstFile, err)
		os.Exit(1)
	}
	defer fp.Close()

	switch templateType {
	case GitIgnore:
		fp.WriteString(templateGitIgnore)
	case ModelBase:
		fp.WriteString(templateModelBase)
	default:
		fmt.Printf("System error: template type %d is not implemented yet\n", templateType)
		os.Exit(1)
	}
}
