package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// cuz in Go "/" is a catch-all
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		homePageHandler(w, r)
	})
	http.HandleFunc("GET /install.sh", bashInstallHandler)
	http.HandleFunc("GET /install.ps1", powershellInstallHandler)

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
		return
	}

	t, err := template.ParseFiles("./web/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func bashInstallHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := os.ReadFile("./scripts/install.sh")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", string(bytes))
}

func powershellInstallHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := os.ReadFile("./scripts/install.ps1")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", string(bytes))
}

func isCurl(r *http.Request) bool {
	userAgent := strings.ToLower(r.UserAgent())
	isCurl := strings.Contains(userAgent, "curl")
	return isCurl
}
