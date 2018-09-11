package mokylin_front

import (
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"strings"
)

var (
	templates *template.Template
)

func initTmpls() {
	files := tmplFiles()
	tmpls, err := parseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}
	templates = tmpls
}

func tmplFiles() []string {
	r := []string{}
	for filename, _ := range _bindata {
		// fmt.Println(filename)
		if strings.HasSuffix(filename, ".html") {
			r = append(r, filename)
		}
	}
	return r
}

// from template.go
func parseFiles(filenames ...string) (*template.Template, error) {
	if len(filenames) == 0 {
		// Not really a problem, but be consistent.
		return nil, fmt.Errorf("html/tmpl: no files named in call to ParseFiles")
	}
	var t *template.Template
	for _, filename := range filenames {
		b, err := Asset(filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		name := filepath.Base(filename)
		// First template becomes return value if not already defined,
		// and we use that one for subsequent New calls to associate
		// all the templates together. Also, if this file has the same name
		// as t, this file becomes the contents of t, so
		//  t, err := New(name).Funcs(xxx).ParseFiles(name)
		// works. Otherwise we create a new template associated with t.
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
