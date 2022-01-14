package main

import (
	"log"
	"net/http"
  "fmt"
  "github.com/devhulk/test-task/runtask"
)

func main() {
    http.HandleFunc("/", runtask.RootHandler)
    http.HandleFunc("/init", runtask.InitHandler)
    port := ":8000"

    fmt.Printf("Run-Task Server running on port %s", port)
    log.Fatal(http.ListenAndServe(port, nil))
}