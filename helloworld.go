package main

import (
    "io"
    "net/http"
    "os"
    "log"
)

func hello(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "Hello world!")
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
  		log.Fatal("$PORT must be set")
  	}

    mux := http.NewServeMux()
    mux.HandleFunc("/", hello)
    http.ListenAndServe(":"+port, mux)
}
