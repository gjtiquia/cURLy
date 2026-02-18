package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", homePageHandler)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3000"
	}
	fmt.Println("running server at port", port)

	err := http.ListenAndServe(":"+port, nil)
	log.Fatal(err)
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}
