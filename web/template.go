package web

import (
	"fmt"
	"os"
)

type templateHelper struct {
	dir string
}

func NewTemplateHelper() *templateHelper {
	return new(templateHelper)
}

func (t *templateHelper) FilePath(templateFile string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("DIR_TEMPLATE"), templateFile)
}

func (t *templateHelper) GetExtendedTemplateFiles(ts string) []string {
	return []string{
		t.FilePath("layout/base.html"),
		t.FilePath(ts),
	}
}
