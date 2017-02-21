package steg

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Config ...
type Config struct {
	StaticRoot  string
	Layout      string
	PartsDir    string
	ContentsDir string
	FuncMap     template.FuncMap
}

// Engine ...
type Engine struct {
	m      map[string]*template.Template
	config Config
}

func parse(dst *template.Template, r io.Reader) error {
	b := new(bytes.Buffer)
	_, err := b.ReadFrom(r)
	if err != nil {
		return err
	}
	_, err = dst.Parse(b.String())
	return err
}

// New ...
func New(config Config) (*Engine, error) {
	base := template.New("")
	dir := http.Dir(config.StaticRoot)
	if err := Walk(dir, config.PartsDir,
		func(path string, info os.FileInfo) error {
			f, err := dir.Open(path)
			if err != nil {
				return fmt.Errorf("read failed in %s: %s", path, err)
			}
			if err := parse(base, f); err != nil {
				return fmt.Errorf("parse failed in %s: %s", path, err)
			}
			return nil
		},
	); err != nil {
		return nil, err
	}
	if config.FuncMap != nil {
		base = base.Funcs(config.FuncMap)
	}
	f, err := dir.Open(config.Layout)
	if err != nil {
		return nil, err
	}
	if err := parse(base, f); err != nil {
		return nil, err
	}
	loads := map[string]*template.Template{}
	if err := Walk(dir, config.ContentsDir,
		func(fname string, info os.FileInfo) error {
			dst, err := base.Clone()
			if err != nil {
				return err
			}
			f, err := dir.Open(fname)
			if err != nil {
				return fmt.Errorf("read failed in %s: %s", fname, err)
			}
			if err := parse(dst, f); err != nil {
				return fmt.Errorf("parse failed in %s: %s", fname, err)
			}
			name, err := filepath.Rel(config.ContentsDir, fname)
			if err != nil {
				return err
			}
			fmt.Println(name)
			loads[name] = dst
			return nil
		},
	); err != nil {
		return nil, err
	}
	return &Engine{
		m:      loads,
		config: config,
	}, nil
}

// ExecuteTemplate ...
func (e *Engine) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	tmpl, ok := e.m[name]
	if !ok {
		return fmt.Errorf("not found template: %q", name)
	}
	return tmpl.Execute(wr, data)
}
