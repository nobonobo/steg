package steg

import (
	"bytes"
	"testing"
)

func TestEngine(t *testing.T) {
	engine, err := New(Config{
		Layout:      "./example/layout.html",
		PartsDir:    "./example/parts",
		ContentsDir: "./example/contents",
		FuncMap:     nil,
	})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	buff := bytes.NewBuffer(nil)
	if err := engine.ExecuteTemplate(buff, "top.html", nil); err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(buff.String())
	buff = bytes.NewBuffer(nil)
	if err := engine.ExecuteTemplate(buff, "blog/index.html", nil); err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(buff.String())
}
