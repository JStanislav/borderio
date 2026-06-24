package gamemanager

import (
	"fmt"
	"time"
)

type EndedGamesCollector struct {
	Threshold int
}

func NewGC(games *Games, threshold int) *EndedGamesCollector {
	fmt.Printf("GC started with threshold %d minute/s\n", threshold)
	ticker := time.NewTicker(time.Duration(threshold) * 60 * time.Second)

	go func() {
		for range ticker.C {
			games.DeleteOldGames()
		}
	}()

	return &EndedGamesCollector{Threshold: threshold}
}
