package irapi

import (
	"context"
	"net/http"
)

type Season struct {
	Active bool `json:"active"`

	Year      int    `json:"year"`
	Quarter   int    `json:"quarter"`
	Week      int    `json:"raceweek"`
	SeasonID  int    `json:"seasonid"`
	SeriesID  int    `json:"seriesid"`
	ShortName string `json:"seriesshortname"`

	LicenceGroup LicenceClass    `json:"serieslicgroupid"`
	Category     LicenceCategory `json:"catid"`

	Start Timestamp `json:"start"`
	End   Timestamp `json:"end"`

	Tracks     []SeasonTrack `json:"tracks"`
	CarClasses []CarClass    `json:"carclasses"`
}

type SeasonTrack struct {
	ID            uint   `json:"id"`
	PackageID     uint   `json:"pkgid"`
	Configuration string `json:"config"`
	Name          string `json:"name"`

	Priority  uint `json:"priority"`
	RaceWeek  int  `json:"raceweek"`
	TimeOfDay int  `json:"timeOfDay"`
}

type CarClass struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	ShortName     string `json:"shortname"`
	RelativeSpeed int    `json:"relspeed"`
	Cars          []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"cars"`
}

type SeasonList []Season

func (c *IRacing) GetSeasons(ctx context.Context, onlyActive bool) (SeasonList, error) {
	path := "/membersite/member/GetSeasons?fields=year,quarter,seriesid,active,catid,carclasses," +
		"tracks,start,end,cars,raceweek,category,serieslicgroupid,names_arr,seasonid,carid" +
		",seriesshortname" +
		"&onlyActive="

	if onlyActive {
		path += "1"
	} else {
		path += "0"
	}

	var seasons SeasonList

	err := c.json(ctx, http.MethodGet, path, nil, &seasons)

	if err != nil {
		return nil, err
	}

	return seasons, nil
}

func (l SeasonList) Len() int { return len(l) }

func (l SeasonList) Less(i, j int) bool {
	a, b := l[i], l[j]

	if a.Year < b.Year {
		return true
	}

	if a.Year == b.Year && a.Quarter < b.Quarter {
		return true
	}

	if a.Quarter == b.Quarter && a.Week < b.Week {
		return true
	}

	return false
}

func (l SeasonList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
