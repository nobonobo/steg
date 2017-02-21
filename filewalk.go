package steg

import (
	"net/http"
	"os"
	"path"
)

// WalkFunc ...
type WalkFunc func(name string, info os.FileInfo) error

// Walk ...
func Walk(d http.Dir, dir string, fn WalkFunc) error {
	return walk(d, dir, fn)
}

func walk(d http.Dir, fname string, fn WalkFunc) error {
	f, err := d.Open(fname)
	if err != nil {
		return err
	}
	v, err := f.Stat()
	if err != nil {
		return err
	}
	if !v.IsDir() {
		if err := fn(fname, v); err != nil {
			return err
		}
		return nil
	}
	files, err := f.Readdir(-1)
	if err != nil {
		return err
	}
	for _, v := range files {
		if err := walk(d, path.Join(fname, v.Name()), fn); err != nil {
			return err
		}
	}
	return nil
}
