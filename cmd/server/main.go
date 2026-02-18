package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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
	bytes, err := os.ReadFile("./markdown/USAGE.md")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isCurl(r) {
		fmt.Fprintf(w, "%s", string(bytes))
	} else {
		// TODO : proper html and styling
		fmt.Fprintf(w, "%s", string(bytes))
	}
}

func isCurl(r *http.Request) bool {
	userAgent := strings.ToLower(r.UserAgent())
	isCurl := strings.Contains(userAgent, "curl")
	return isCurl
}
