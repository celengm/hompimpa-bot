package main

import (
    "io"
    "net/http"
    "fmt"
    //"os"
    //"log"
)

func hello(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if len(name) != 0 {
      io.WriteString(w, "Hello " + name)
    }
}

func callback() {
    defer req.Body.Close()
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
      fmt.Println(err)
    }
    decoded, err := base64.StdEncoding.DecodeString(req.Header.Get("X-Line-Signature"))
    if err != nil {
      fmt.Println(err)
    }
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
