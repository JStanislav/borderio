package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JStanislav/quoridor-clone/config"
	"github.com/JStanislav/quoridor-clone/external"
	"github.com/JStanislav/quoridor-clone/gamemanager"
	ws "github.com/JStanislav/quoridor-clone/websocket"
)

func main() {
	config := config.LoadConfig()

	localhost := "127.0.0.1"
	fmt.Printf("Server is running on %s:%s\n", localhost, config.Port)

	mux := http.NewServeMux()

	games := gamemanager.NewGames()
	updateStatsServiceClient := external.NewUpdateStatsServiceHTTPClient("") // a NATS client can be used too, just for fun.

	wsHandler := ws.NewHandler(&games, updateStatsServiceClient)

	mux.HandleFunc("/{id}", wsHandler.Handler)
	mux.HandleFunc("/ping/{hash}", wsHandler.GamePing)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", localhost, config.Port), mux))
}
