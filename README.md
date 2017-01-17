# steg

Simple Template Engine for Go

# feature

- layout template
- part define
- content definitions
- add global func-map

# install

```sh
go get -u github.com/nobonobo/steg
```

# use

```go
package main

import "github.com/nobonobo/steg"

func main() {
    engine, err := steg.New(steg.Config{
        Layout:   "./example/layout.html",
        PartsDir: "./example/parts",
        ContentsDir: "./example/contents",
    })
    if err != nil {
        // error handling
    }
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
        if err := engine.ExecuteTemplate(w, "top.html", nil); err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
    })
    ...
}
```

