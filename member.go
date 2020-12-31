package irapi

import (
	"context"
	"net/http"
	"strconv"
)

// CareerStats represents a member's lifetmie career stats for a single category
type CareerStats struct {
	Wins                    uint    `json:"wins"`
	TotalClubPoints         uint    `json:"totalclubpoints"`
	WinPercentage           float64 `json:"winperc"`
	Poles                   uint    `json:"poles"`
	AverageStart            uint    `json:"avgStart"`
	AverageFinish           uint    `json:"avgFinish"`
	TopFivePercent          float64 `json:"top5Perc"`
	TotalLaps               uint64  `json:"totalLaps"`
	AverageIncidentsPerRace float64 `json:"avgIncPerRace"`
	AveragePointsPerRace    float64 `json:"avgPtsPerRace"`
	LapsLed                 uint    `json:"lapsLed"`
	TopFiveFinishes         uint    `json:"top5"`
	LapsLedPercentage       float64 `json:"lapsLedPerc"`
	Category                string  `json:"category"`
	Starts                  uint    `json:"starts"`
}

// GetCareerStats gets the lifetime career stats for a user
func (c *IRacing) GetCareerStats(ctx context.Context, userID uint64) ([]CareerStats, error) {
	path := "/memberstats/member/GetCareerStats?custid=" + strconv.FormatUint(userID, 10)

	careerStats := []CareerStats{}

	err := c.json(ctx, http.MethodGet, path, nil, &careerStats)

	if err != nil {
		return nil, err
	}

	return careerStats, nil
}
