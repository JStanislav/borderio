package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func handle(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("[ERROR] error upgrading, %s\n", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			fmt.Printf("[ERROR] error reading message, %s\n", err)
			break
		}

		fmt.Printf("received %s\n", message)
		err = c.WriteMessage(mt, []byte("pong"))
		if err != nil {
			fmt.Printf("[ERROR] error writing message, %s\n", err)
			break
		}
	}
}

func main() {
	fmt.Println("Hello World")

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
