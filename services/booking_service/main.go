package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialize the service
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	log.Println("Server running on :8002")
	log.Fatal(http.ListenAndServe(":8002", nil))
}
