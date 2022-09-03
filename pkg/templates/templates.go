package templates

const (
	TemplateGitIgnore int = iota
)

func Exec(templateType int, dstFile string, data any) {
	switch templateType {
	case TemplateGitIgnore:
	default:
	}
}
