package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JStanislav/quoridor-clone/config"
	"github.com/JStanislav/quoridor-clone/external"
	"github.com/JStanislav/quoridor-clone/gamemanager"
	middleware "github.com/JStanislav/quoridor-clone/middlewares"
	ws "github.com/JStanislav/quoridor-clone/websocket"
)

func main() {
	config := config.LoadConfig()

	localhost := "0.0.0.0"
	fmt.Printf("Server is running on %s:%s\n", localhost, config.Port)

	mux := http.NewServeMux()

	gamesContainer := gamemanager.NewGamesContainer(1)
	updateStatsServiceClient := external.NewUpdateStatsServiceHTTPClient("") // a NATS client can be used too, just for fun.

	handlerContext := context.WithValue(context.Background(), "TimeoutAfterGameOver", time.Duration(config.TimeoutAfterGameOver)*time.Second)

	games := gamemanager.Games(gamesContainer.Games)
	wsHandler := ws.NewHandler(handlerContext, &games, updateStatsServiceClient)

	mux.HandleFunc("/{id}", wsHandler.Handler)
	mux.HandleFunc("/ping/{hash}", wsHandler.GamePing)
	mux.HandleFunc("/game_stats", wsHandler.GamesList)

	configMiddleware := middleware.NewCORSConfig(config.Cors.AllowedOrigins, config.Cors.AllowedMethods)
	handler := middleware.CORS(configMiddleware)(mux)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", localhost, config.Port), handler))
}
