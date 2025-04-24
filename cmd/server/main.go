package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting API server...")
	// TODO: Implement API endpoints

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to NSight Proxy!")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
