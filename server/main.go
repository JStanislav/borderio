package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JStanislav/quoridor-clone/config"
	"github.com/JStanislav/quoridor-clone/game"
	ws "github.com/JStanislav/quoridor-clone/websocket"
)

var games = make(map[string]*game.GameState)

func createHash(h string) *game.GameState {
	games[h] = nil
	return games[h]
}

func getGame(h string) *game.GameState {
	return games[h]
}

func main() {
	config := config.LoadConfig()

	localhost := "127.0.0.1"
	fmt.Printf("Server is running on %s:%s\n", localhost, config.Port)

	mux := http.NewServeMux()

	wsHandler := ws.NewHandler(createHash, getGame)

	mux.HandleFunc("/{id}", wsHandler.Handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", localhost, config.Port), mux))
}
