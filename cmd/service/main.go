package main

import (
	"fmt"
	"net/http"
	"os"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ANWORK")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	r.Write(w)
}

func main() {
	fmt.Println("Hello!")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/echo", echoHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil); err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}
}
