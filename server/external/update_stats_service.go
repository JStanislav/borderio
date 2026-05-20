package external

import (
	"fmt"

	"github.com/JStanislav/quoridor-clone/game"
	"github.com/go-resty/resty/v2"
)

type UpdateStats func(GameStats game.GameStats) error

type UpdateStatsService interface {
	UpdateStats(GameStats game.GameStats) error
}

type UpdateStatsServiceHTTPClient struct {
	restyClient *resty.Client
	ServiceURL  string
}

func NewUpdateStatsServiceHTTPClient(serviceURL string) *UpdateStatsServiceHTTPClient {
	return &UpdateStatsServiceHTTPClient{
		restyClient: resty.New(),
		ServiceURL:  serviceURL,
	}
}

func (c *UpdateStatsServiceHTTPClient) UpdateStats(gameStats game.GameStats) error {
	fmt.Println("Function yet not implemented.")
	return nil
}
