package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Tunnel hit! Path: %s | User-Agent: %s\n", r.URL.Path, r.UserAgent())
		fmt.Fprintf(w, "<h1>It works!</h1><p>Served from my local laptop via Azure ACI (4.187.182.59).</p>")
	})

	log.Println("Local HTTP server listening on 127.0.0.1:6767...")
	if err := http.ListenAndServe("127.0.0.1:6767", nil); err != nil {
		log.Fatal(err)
	}
}
