package main

import (
  "fmt"
  "net/http"
  "io"
  "log"
)

func main() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  	fmt.Println("received request")

    b := make([]byte, 8)
    for {
      n, err := r.Body.Read(b)
      fmt.Printf("read %d bytes: %s\n", n, string(b))
      if err == io.EOF {
        fmt.Println("reached EOF")
        break
      }
    }

    defer func() {
      fmt.Println("done")
    }()
  })

fmt.Println("server listening on port 8080")

log.Fatal(http.ListenAndServe(":8080", nil))
}
