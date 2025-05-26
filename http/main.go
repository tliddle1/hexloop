package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./static")) // Serve files from the "static" folder
	http.Handle("/", fs)

	port := "8081"
	log.Println("Serving on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
