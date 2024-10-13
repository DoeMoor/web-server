package main

import (
	// "fmt"
	"log"
	"net/http"
	"time"
)



func main() {
	
	mux := http.NewServeMux()
	serverStruct := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		ReadTimeout: 2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout: 2 * time.Second,
	}

	mux.Handle("GET /", http.FileServer(http.Dir(".")))
	mux.Handle("GET /assets/", http.FileServer(http.Dir(".")))
	
	log.Printf("Serving on port: %s\n", serverStruct.Addr)
	log.Fatal(serverStruct.ListenAndServe())

}
