package main

import (
	"fmt"
	"log"
	"net/http"

	ws "github.com/JStanislav/quoridor-clone/websocket"
)

func main() {
	fmt.Println("Server is running on localhost:8080")

	http.HandleFunc("/", ws.Handle)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
