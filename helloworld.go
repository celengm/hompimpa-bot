package main

import (
    "io"
    "io/ioutil"
    "net/http"
    "fmt"
    "os"
    "log"
    "encoding/base64"
    "crypto/hmac"
    "crypto/sha256"
)

func hello(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if len(name) != 0 {
      io.WriteString(w, "Hello " + name)
    }
}

func callback(w http.ResponseWriter, req *http.Request) {
    defer req.Body.Close()
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
      fmt.Println(err)
    }
    fmt.Println(body)
    decoded, err := base64.StdEncoding.DecodeString(req.Header.Get("X-Line-Signature"))
    if err != nil {
      fmt.Println(err)
    }
    fmt.Println(decoded)
    hash := hmac.New(sha256.New, []byte("6db72166ed2b37fbfd0a4a00f7bd01ac"))
    hash.Write(body)
}

func main() {

    port := os.Getenv("PORT")
    if port == "" {
  		log.Fatal("$PORT must be set")
  	}


    mux := http.NewServeMux()
    mux.HandleFunc("/", hello)
    mux.HandleFunc("/callback", callback)
    // http.ListenAndServe(":8000", mux)
    http.ListenAndServe(":"+port, mux)
}
