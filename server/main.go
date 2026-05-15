package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JStanislav/quoridor-clone/config"
	"github.com/JStanislav/quoridor-clone/gamemanager"
	ws "github.com/JStanislav/quoridor-clone/websocket"
)

func main() {
	config := config.LoadConfig()

	localhost := "127.0.0.1"
	fmt.Printf("Server is running on %s:%s\n", localhost, config.Port)

	mux := http.NewServeMux()

	games := gamemanager.NewGames()

	wsHandler := ws.NewHandler(&games)

	mux.HandleFunc("/{id}", wsHandler.Handler)
	mux.HandleFunc("/ping/{hash}", wsHandler.GamePing)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", localhost, config.Port), mux))
}
