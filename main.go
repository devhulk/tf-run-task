package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/devhulk/test-task/runtask"
)

func main() {
	http.HandleFunc("/", runtask.TaskHandler)
	port := ":8000"

	fmt.Printf("Run-Task Server running on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
