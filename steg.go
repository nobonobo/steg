package steg

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Config ...
type Config struct {
	Layout      string
	PartsDir    string
	ContentsDir string
	FuncMap     template.FuncMap
}

// Engine ...
type Engine struct {
	config Config
	m      map[string]*template.Template
}

// New ...
func New(config Config) (*Engine, error) {
	parts, err := template.New("").ParseGlob(filepath.Join(config.PartsDir, "*.html"))
	if err != nil {
		return nil, err
	}
	if config.FuncMap != nil {
		parts = parts.Funcs(config.FuncMap)
	}
	b, err := ioutil.ReadFile(config.Layout)
	if err != nil {
		return nil, err
	}
	layout, err := parts.Parse(string(b))
	if err != nil {
		return nil, err
	}
	m := map[string]*template.Template{}
	if err := filepath.Walk(
		config.ContentsDir,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				name, err := filepath.Rel(config.ContentsDir, path)
				if err != nil {
					return err
				}
				b, err := ioutil.ReadFile(path)
				if err != nil {
					return fmt.Errorf("parse fail in %s: %s", name, err)
				}
				clone, err := layout.Clone()
				if err != nil {
					return fmt.Errorf("parse fail in %s: %s", name, err)
				}
				parsed, err := clone.Parse(string(b))
				if err != nil {
					return fmt.Errorf("parse fail in %s: %s", name, err)
				}
				m[name] = parsed
			}
			return nil
		},
	); err != nil {
		return nil, err
	}
	return &Engine{
		config: config,
		m:      m,
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
