package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/stianeikeland/go-rpio/v4"
)

const ledPin = 26

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello!")

}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	pin := rpio.Pin(ledPin)
	pin.Toggle()
	fmt.Fprint(w, "Led toggled")
}

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}
	defer rpio.Close()
	http.HandleFunc("/my", helloHandler)
	http.HandleFunc("/toggle", toggleHandler)
	log.Println("Starting server")
	http.ListenAndServe(":7070", nil)
}
