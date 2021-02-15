package web

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type templateHelper struct {
	dir string
}

func NewTemplateHelper() *templateHelper {
	return new(templateHelper)
}

func (t *templateHelper) Dir() string {
	d, e := filepath.Abs(filepath.Dir(os.Args[0]))
	if e != nil {
		log.Fatal(e)
	}
	return d
}

func (t *templateHelper) FilePath(templateFile string) string {
	return fmt.Sprintf("%s/template/%s", t.Dir(), templateFile)
}

func (t *templateHelper) GetExtendedTemplateFiles(ts string) []string {
	return []string{
		t.FilePath("layout/base.html"),
		t.FilePath(ts),
	}
}
