# itunesService
service to consulting itunes API

# Install
```
go get github.com/serter95/itunesServiceGo
```
# Use

```go
package main

import "itunesServiceGo"

func main() {
    sliceWithData, err := itunesServiceGo.FindResults("your criteria")
    // do wath you want
}

// the slice data will content this struct:

type StandardResponse struct {
	Category   string `json:"category"`
	Name       string `json:"name"`
	Author     string `json:"author"`
	PreviewURL string `json:"previewUrl"`
	Origin     string `json:"origin"`
}
```