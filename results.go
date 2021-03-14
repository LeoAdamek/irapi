package irapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// SortField is an enum representing the fields which can be used to sort results
type SortField uint8

// SortDir is an enum of the directions a set of results can be sorted by (ASC and DESC)
type SortDir uint8

const (
	SortFieldStartTime SortField = iota
	SortFieldSeasonQuarter
	SortFieldRaceWeek
	SortFieldEventType
	SortFieldSeriesName
	SortFieldClassAndCar
	SortFieldTrackName
	SortFieldStartPos
	SortFieldFinishPos
	SortFieldIncidents
	SortFieldChampionshipPoints
	SortFieldSOF
	SortFieldWinnerName

	SortDirASC SortDir = iota
	SortDirDESC

	SessionPhaseRace = iota
)

var resultSortFieldStrings = []string{
	"start_time",
	"season_quarter",
	"race_week_num",
	"evttypename",
	"series_shortname",
	"car_class_name",
	"track_name",
	"starting_position",
	"finishing_position",
	"incidents",
	"champpoints",
	"strengthoffield",
	"winnerdisplayname",
}

func (f SortField) String() string {
	if f > SortFieldWinnerName {
		return fmt.Sprintf("Invalid SortField: %d", f)
	}

	return resultSortFieldStrings[int(f)]
}

func (d SortDir) String() string {
	switch d {
	case SortDirDESC:
		return "desc"
	case SortDirASC:
		return "asc"
	default:
		panic(fmt.Sprintf("SortOrder out of range: %d", d))
	}
}

// SessionResult is the result data of a single session
type SessionResult struct {
	ID uint64 `json:"subsessionid"`

	AverageLapTime          Milliseconds    `json:"eventavglap"`
	CategoryID              LicenceCategory `json:"catid"`
	CautionLaps             uint            `json:"cautionlaps"`
	CautionType             int8            `json:"cautiontype"`
	Cautions                uint            `json:"cautions"`
	CornersPerLap           uint8           `json:"cornersperlap"`
	DriverChangeRule        int8            `json:"driver_change_rule"`
	DriverChanges           uint            `json:"driver_changes"`
	EventType               int8            `json:"evttype"`
	LapsComplete            uint            `json:"eventlapscompleted"`
	LapsForSoloAverage      uint            `json:"nlapsforsoloavg"`
	LeadChanges             int             `json:"nleadchanges"`
	LeaveMarbles            uint8           `json:"leavemarbles"`
	MaxWeeks                uint8           `json:"maxweeks"`
	MaximumTeamDrivers      uint8           `json:"max_team_drivers"`
	MinimumTeamDrivers      uint8           `json:"min_team_drivers"`
	PointsType              string          `json:"pointstype"`
	PrivateSessionID        int64           `json:"privatesessionid"`
	Quarter                 int8            `json:"season_quarter"`
	RaceWeek                uint            `json:"race_week_num"`
	RubberLevelPractice     int8            `json:"rubberlevel_practice"`
	RubberLevelQualify      int8            `json:"rubberlevel_qualify"`
	RubberLevelRace         int8            `json:"rubberlevel_race"`
	RubberLevelWarmup       int8            `json:"rubberlevel_warmup"`
	SOF                     uint            `json:"eventstrengthoffield"`
	SeasionShortname        string          `json:"seasion_shortname"`
	SeasonID                int64           `json:"seasonID"`
	SeasonName              string          `json:"season_name"`
	SeasonYear              uint            `json:"season_year"`
	SeriesID                uint64          `json:"seriesid"`
	SeriesName              string          `json:"series_name"`
	SessionID               uint64          `json:"sessionid"`
	SessionName             string          `json:"sessionname"`
	SimSessionType          int8            `json:"simsestype"`
	SimulatedStartTime      SimTime         `json:"simulatedstarttime"`
	SpecialEventType        int8            `json:"specialeventtype"`
	StartTime               SimTime         `json:"starttime"`
	TimeOfDay               int8            `json:"timeofday"`
	TrackConfigName         string          `json:"track_config_name"`
	TrackID                 uint64          `json:"trackid"`
	TrackName               string          `json:"track_name"`
	WeatherFogDensity       float64         `json:"weather_fog_density"`
	WeatherHumidity         float64         `json:"weather_rh"`
	WeatherSkies            int8            `json:"weather_skies"`
	WeatherTemperatureUnit  int8            `json:"weather_temp_units"`
	WeatherTemperatureValue float64         `json:"weather_temp_value"`
	WeatherType             int8            `json:"weather_type"`
	WeatherVariation        int             `json:"weather_var_ongoing"`
	WeatherWindDirection    int8            `json:"weather_wind_dir"`
	WeatherWindspeedUnits   int8            `json:"weather_wind_speed_units"`
	WeatherWindspeedValue   float64         `json:"weather_wind_speed_value"`

	Results []DriverResult `json:"rows"`
}

type getLapTimesResponse struct {
	Laptimes []LapResult `json:"lapData"`
}

type LapFlags uint64

type LapResult struct {
	SessionTime uint64   `json:"ses_time"`
	UserID      uint64   `json:"custid"`
	Flags       LapFlags `json:"flags"`
	LapNumber   uint64   `json:"lap_num"`
}

// DriverResult shows the results for a single driver
type DriverResult struct {
	Name           string `json:"displayname"`
	CarNumber      string `json:"carnum"`
	StartPosition  int    `json:"startpos"`
	FinishPosition int    `json:"finishpos"`

	BestNLapNumber int     `json:"bestlapnum"`
	NewCPI         float64 `json:"newcpi"`
	SessionName    string  `json:"simsesname"`
	CarClassName   string  `json:"ccName"`
	OldIRating     int     `json:"oldirating"`
	NewIRating     int     `json:"newirating"`
	CarID          uint    `json:"carid"`
	LapsCompleted  uint    `json:"lapscomplete"`
}

// SearchResultsOptions represents various options which can be given when searching results
type SearchResultsOptions struct {
	IncludeRaces          bool
	IncludeQualifications bool
	IncludeTimeTrials     bool
	IncludeOPs            bool
	IncludeOfficial       bool
	IncludeUnofficial     bool
	IncludeRookie         bool
	IncludeClassD         bool
	IncludeClassC         bool
	IncludeClassB         bool
	IncludeClassA         bool
	IncludePro            bool
	IncludeProWC          bool

	Season *SeasonFilter

	DateRange *DateRange
	UserID    uint64

	SortBy  SortField
	SortDir SortDir
}

type SeasonFilter struct {
	SeasonYear    int
	SeasonQuarter int
	SeasonWeek    *int
}

type DateRange struct {
	Lower time.Time
	Upper time.Time
}

// DefaultSearchResultsOptions gets the default set of options for searching results
//
// The default options will show only and all official races in all classes for the current season
func DefaultSearchResultsOptions() *SearchResultsOptions {
	r := &SearchResultsOptions{
		IncludeRaces:          true,
		IncludeQualifications: false,
		IncludeOPs:            false,
		IncludeTimeTrials:     false,
		IncludeOfficial:       true,
		IncludeUnofficial:     false,
		IncludeRookie:         true,
		IncludeClassD:         true,
		IncludeClassC:         true,
		IncludeClassB:         true,
		IncludeClassA:         true,
		IncludePro:            true,
		IncludeProWC:          true,

		SortBy:  SortFieldStartTime,
		SortDir: SortDirDESC,
	}

	return r
}

type searchResultsResponse struct {
	Headers map[string]string `json:"m"`
	Data    struct {
		Count int                `json:"15"`
		Rows  []SearchResultData `json:"r"`
	} `json:"d"`
}

// SearchResultData represents the data results for a search
type SearchResultData struct {
	HelmetColor1            string    `json:"1"`
	WinnerHelmetColor2      string    `json:"2"`
	FinishPos               int       `json:"3"`
	WinnerHelmetColor3      string    `json:"4"`
	WinnerHelmetColor4      string    `json:"5"`
	BestQualifictionLapTime string    `json:"6"`
	SubSessionBestLapTime   string    `json:"7"`
	RaceWeek                int       `json:"8"`
	SessionID               uint64    `json:"9"`
	FinishedAt              Timestamp `json:"10"`
	RawStartTime            Timestamp `json:"11"`
	StartingPos             int       `json:"12"`
	HelmetColor3            string    `json:"13"`
	HelmetColor2            string    `json:"14"`
	//RowCount                int          `json:"15"`
	ClubPoints             int             `json:"16"`
	DropRacePoints         int             `json:"17"`
	OfficialSession        int             `json:"18"`
	GroupName              string          `json:"19"`
	SeriesID               int             `json:"20"`
	StartTime              string          `json:"21"`
	SeasonID               int             `json:"22"`
	UserID                 uint64          `json:"23"`
	HelmetLicenceLevel     LicenceClass    `json:"24"`
	WinnerLicenseLevel     LicenceClass    `json:"25"`
	RowNumber              int             `json:"26"`
	WinnersGroupID         int             `json:"27"`
	SessionRank            int             `json:"28"`
	CarClassID             uint64          `json:"29"`
	TrackID                uint64          `json:"30"`
	WinnerName             string          `json:"31"`
	CarID                  uint64          `json:"32"`
	CategoryID             LicenceCategory `json:"33"`
	SeasonQuarter          int8            `json:"34"`
	LicenceGroup           LicenceClass    `json:"35"`
	WinnerHelmetPattern    int             `json:"36"`
	EventType              int             `json:"37"`
	BestLapTime            string          `json:"38"`
	Incidents              int             `json:"39"`
	ChapionshipPoints      int             `json:"40"`
	SubsessionID           uint64          `json:"41"`
	SeasonYear             int             `json:"42"`
	ChampionshipPointsSort int             `json:"43"`
	StartDate              string          `json:"44"`
	SOF                    int             `json:"45"`
	HelmetPattern          int             `json:"46"`
	ClubPointsSort         int             `json:"47"`
	DisplayName            string          `json:"48"`
}

// GetSubSessionResult gets the result of a single iRacing subsession (often referred to as a "split")
func (c *IRacing) GetSubSessionResult(ctx context.Context, subsessionID uint64) (*SessionResult, error) {

	path := "/membersite/member/GetSubsessionResults?subsessionID=" + strconv.FormatUint(subsessionID, 10)
	result := &SessionResult{}

	if err := c.json(ctx, http.MethodGet, path, nil, result); err != nil {
		return nil, err
	}

	return result, nil
}

// SearchResults searches for results based on the ID of a given participent user and other options
//
// @param userID ID of a user who participated in the session
func (c *IRacing) SearchResults(ctx context.Context, opts *SearchResultsOptions) ([]SearchResultData, error) {

	if opts.DateRange != nil && opts.Season != nil {
		return nil, errors.New("only one of Season or DateRange may be specified")
	} else if opts.DateRange == nil && opts.Season == nil {
		return nil, errors.New("one of DateRange or Season must be specified")
	}

	path := "/memberstats/member/GetResults"
	params := make(url.Values)

	if opts.UserID > 0 {
		params.Set("custid", strconv.FormatUint(opts.UserID, 10))
	}

	if d := opts.DateRange; d != nil {
		params.Set("starttime_low", strconv.FormatInt(d.Lower.Unix()*1000, 10))
		params.Set("starttime_high", strconv.FormatInt(d.Upper.Unix()*1000, 10))
	}

	params.Set("lowerbound", strconv.FormatInt(0, 10))
	params.Set("upperbound", strconv.FormatInt(100, 10))

	fm := func(f bool) string {
		i := uint64(0)

		if f {
			i = 1
		}

		return strconv.FormatUint(i, 10)
	}

	params.Set("showraces", fm(opts.IncludeRaces))
	params.Set("showquals", fm(opts.IncludeQualifications))
	params.Set("showtts", fm(opts.IncludeTimeTrials))
	params.Set("showofficial", fm(opts.IncludeOfficial))
	params.Set("showunofficial", fm(opts.IncludeUnofficial))
	params.Set("showrookie", fm(opts.IncludeRookie))
	params.Set("showclassd", fm(opts.IncludeClassD))
	params.Set("showclassc", fm(opts.IncludeClassC))
	params.Set("showclassb", fm(opts.IncludeClassB))
	params.Set("showclassa", fm(opts.IncludeClassA))
	params.Set("showpro", fm(opts.IncludePro))
	params.Set("showprowc", fm(opts.IncludeProWC))

	params.Set("category[]", "1,2,3,4")
	params.Set("format", "json")
	params.Set("sort", opts.SortBy.String())
	params.Set("order", opts.SortDir.String())

	path += "?" + params.Encode()

	resp := &searchResultsResponse{}

	if err := c.json(ctx, http.MethodGet, path, nil, resp); err != nil {
		return nil, err
	}

	return resp.Data.Rows, nil
}

func (c *IRacing) GetLaps(ctx context.Context, sessionID uint64, entrantID uint64, phase uint64) ([]LapResult, error) {
	path := "/membersite/member/GetLaps"

	params := make(url.Values)
	params.Set("subsessionid", strconv.FormatUint(sessionID, 10))
	params.Set("groupid", strconv.FormatUint(entrantID, 10))
	params.Set("simsesnum", strconv.FormatUint(phase, 10))

	path += "?=&" + params.Encode()

	resp := &getLapTimesResponse{}

	if err := c.json(ctx, http.MethodPost, path, strings.NewReader("a=null"), resp); err != nil {
		return nil, err
	}

	return resp.Laptimes, nil
}
