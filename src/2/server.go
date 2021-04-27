package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Listening on port :8000\nlocalhost:8000")
	log.Fatal(http.ListenAndServe(":8000", http.FileServer(http.Dir("./"))))
}
