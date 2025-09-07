package datasource

import (
	"context"
	"io"

	"fmt"
	"maps"
	"net/http"
	"time"

	"cosmossdk.io/log"
	"github.com/goccy/go-json"
)

// It is demonstration datasource for futchain module.
// This module is responsible to feed the chain with football data.
// This module works as oracle for the chain.
// Aiming to have consensus over football data in future. For now, we only have 1 data source.

type DatasourceFM struct {
	Client  *http.Client // will apply default h2 optimizations ,need stealth client?
	BaseURL string
	Headers http.Header
}

type FetchSettings func(f *fetchParams)

func WithTimezone(timezone string) FetchSettings {
	return FetchSettings(func(f *fetchParams) {
		f.timezone = timezone
	})
}

func WithLogger(logger log.Logger) FetchSettings {
	return FetchSettings(func(f *fetchParams) {
		f.logger = logger
	})
}

type fetchParams struct {
	timezone string
	logger   log.Logger
}

func (d *DatasourceFM) Fetch(ctx context.Context, s ...FetchSettings) ([]League, error) {

	var params fetchParams

	for _, s := range s {
		s(&params)
	}

	// Get current time in the specified timezone
	loc, err := time.LoadLocation(params.timezone)
	if err != nil {
		return nil, err
	}
	now := time.Now().In(loc)
	date := now.Format("20060102")

	//https://www.fotmob.com/api/data/matches?date=20250906&timezone=Europe%2FIstanbul&ccode3=GBR

	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/data/matches?date=%s&timezone=%s&ccode3=GBR", d.BaseURL, date, params.timezone), nil)
	if err != nil {
		return nil, err
	}
	request.Header = maps.Clone(d.Headers)

	response, err := d.Client.Do(request)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if params.logger != nil {
			params.logger.Error("error fetching data", "status", response.StatusCode)
		}

		return nil, fmt.Errorf("error fetching data: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body) //TODO: optimize it with sync.pool
	if err != nil {
		if params.logger != nil {
			params.logger.Error("error reading response body", "error", err)
		}

		return nil, err
	}

	var result struct {
		League []League `json:"leagues"`
	}

	return result.League, json.Unmarshal(body, &result)

}

type League struct {
	IsGroup   bool    `json:"isGroup"`
	GroupName string  `json:"groupName"`
	Ccode     string  `json:"ccode"`
	ID        int     `json:"id"`
	PrimaryID int     `json:"primaryId"`
	Name      string  `json:"name"`
	Matches   []Match `json:"matches"`
}

type Match struct {
	ID               int    `json:"id"`
	LeagueID         int    `json:"leagueId"`
	Time             string `json:"time"`
	Home             Team   `json:"home"`
	Away             Team   `json:"away"`
	EliminatedTeamID any    `json:"eliminatedTeamId"`
	StatusID         int    `json:"statusId"`
	TournamentStage  string `json:"tournamentStage"`
	Status           Status `json:"status"`
	Ongoing          bool   `json:"ongoing"`
	TimeTS           int64  `json:"timeTS"`
}
type Team struct {
	ID       int    `json:"id"`
	Score    int    `json:"score"`
	Name     string `json:"name"`
	LongName string `json:"longName"`
}

type Halfs struct {
}
type Status struct {
	UtcTime      time.Time `json:"utcTime"`
	Halfs        Halfs     `json:"halfs"`
	PeriodLength int       `json:"periodLength"`
	Started      bool      `json:"started"`
	Cancelled    bool      `json:"cancelled"`
	Finished     bool      `json:"finished"`
	Ongoing      bool      `json:"ongoing"`
	LiveTime     LiveTime  `json:"liveTime"`
}

type LiveTime struct {
	Long      string `json:"long"`      // "long": "51:35",
	MaxTime   int    `json:"maxTime"`   // "maxTime": 90,
	AddedTime int    `json:"addedTime"` // "addedTime": 0
}

type ComparePriority int

func (c ComparePriority) EventName() string {
	return priorityNamer[c]
}

const (
	PriorityNoChanges ComparePriority = iota
	PriorityLiveTime
	PriorityPeriodLength
	PriorityStatus
	PriorityOngoing
	PriorityFinished
	PriorityStarted
	PriorityCancelled
	PriorityScore
)

var priorityNamer = []string{
	"match_no_changes",
	"match_live_time",
	"match_period_length",
	"match_status",
	"match_ongoing",
	"match_finished",
	"match_started",
	"match_cancelled",
	"match_score",
}

const MinimumEventPriority ComparePriority = PriorityPeriodLength

func (new *Match) Compare(old *Match) ComparePriority {
	// check changes in reverse priority order: the higher the priority, the more important the change is
	switch {
	case new.Home.Score != old.Home.Score:
		return PriorityScore
	case new.Away.Score != old.Away.Score:
		return PriorityScore
	case new.Status.Cancelled != old.Status.Cancelled:
		return PriorityCancelled
	case new.Status.Finished != old.Status.Finished:
		return PriorityFinished
	case new.Status.Started != old.Status.Started:
		return PriorityStarted
	case new.Status.Ongoing != old.Status.Ongoing:
		return PriorityOngoing
	case new.Status.PeriodLength != old.Status.PeriodLength:
		return PriorityPeriodLength
	case new.Status.LiveTime.Long != old.Status.LiveTime.Long:
		return PriorityLiveTime
	case new.Status.LiveTime.MaxTime != old.Status.LiveTime.MaxTime:
		return PriorityLiveTime
	case new.Status.LiveTime.AddedTime != old.Status.LiveTime.AddedTime:
		return PriorityLiveTime
	}
	return PriorityNoChanges
}
