package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JStanislav/quoridor-clone/config"
	ws "github.com/JStanislav/quoridor-clone/websocket"
)

func main() {
	config := config.LoadConfig()

	localhost := "127.0.0.1"
	fmt.Printf("Server is running on %s:%s\n", localhost, config.Port)

	http.HandleFunc("/", ws.Handle)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", localhost, config.Port), nil))
}
