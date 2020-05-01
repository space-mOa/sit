package main

import (
	"log"
	"net/http"
	"fmt"
)

func main() {
	fmt.Println("Listening on port :8000")
	log.Fatal(http.ListenAndServe(":8000", http.FileServer(http.Dir("./"))))
}
