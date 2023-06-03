# req

A request tool based on Go native `net/http`.

## Getting `req`

With Go module support, simply add the following import

```go
import "github.com/iamdanielyin/req"
```
to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, run the following Go command to install the gin package:

```shell
$ go get -u github.com/iamdanielyin/req
```

## Running `req`

First you need to import `req` package for using `req`, one simplest example likes the follow example.go:

```go
package main

import (
	"github.com/iamdanielyin/req"
	"log"
)

func main() {
	var posts []struct {
		UserId int    `json:"userId"`
		Id     int    `json:"id"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}
	if err := req.GET("https://jsonplaceholder.typicode.com/posts", &posts); err != nil {
		log.Fatal(err)
	}
}
```
And use the Go command to run the demo:
```go
$ go run example.go
```

## List of APIs

- `req.GET(url, dst)`
- `req.POST(url, body, dst)`
- `req.PATCH(url, body, dst)`
- `req.PUT(url, body, dst)`
- `req.DELETE(url, body, dst)`
- `req.CALL(method, url, body, dst)`
- `req.URL(format, values...)`
  - `.GET(dst[, headers])`
  - `.POST(body, dst[, headers])`
  - `.PATCH(body, dst[, headers])`
  - `.PUT(body, dst[, headers])`
  - `.DELETE(dst[, headers])`
  - `.CALL(method, body, dst[, headers])`
  - `.AddHeader(key, value)`
  - `.DelHeader(key)`
  - `.SetHeaders(headers)`
  - `.SetBody(body)`
  - `.URL()`
  - `.Headers()`