package gamemanager

import (
	"fmt"
	"log/slog"
	"time"
)

type EndedGamesCollector struct {
	Threshold int
}

func NewGC(games *Games, threshold int) *EndedGamesCollector {
	thresholdStr := fmt.Sprintf("%d minute/s", threshold)
	slog.Info("GC Started", "threshold", thresholdStr)

	ticker := time.NewTicker(time.Duration(threshold) * 60 * time.Second)

	go func() {
		for range ticker.C {
			games.DeleteOldGames()
		}
	}()

	return &EndedGamesCollector{Threshold: threshold}
}
