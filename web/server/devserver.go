package main

import (
	"log"
	"net/http"
)

func main() {
	public := http.Dir("./web/public")

	devserver := http.FileServer(public)
	http.Handle("/", devserver)

	log.Println("Server listening port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
