package flatbuffers

import (
	"errors"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/raifpy/futchain/x/futchain/keeper/datasource"
	"github.com/raifpy/futchain/x/futchain/keeper/datasource/flatbuffers/futchain"
)

// EncodeTeam encodes a Team struct to FlatBuffers binary format
func EncodeTeam(team *datasource.Team) ([]byte, error) {
	if team == nil {
		return nil, errors.New("team cannot be nil")
	}

	builder := flatbuffers.NewBuilder(1024)

	// Create string offsets
	nameOffset := builder.CreateString(team.Name)
	longNameOffset := builder.CreateString(team.LongName)

	// Create the Team table
	futchain.TeamStart(builder)
	futchain.TeamAddId(builder, int32(team.ID))
	futchain.TeamAddScore(builder, int32(team.Score))
	futchain.TeamAddName(builder, nameOffset)
	futchain.TeamAddLongName(builder, longNameOffset)
	teamOffset := futchain.TeamEnd(builder)

	builder.Finish(teamOffset)
	return builder.FinishedBytes(), nil
}

// DecodeTeam decodes FlatBuffers binary data to a Team struct
func DecodeTeam(data []byte) (*datasource.Team, error) {
	if len(data) == 0 {
		return nil, errors.New("data cannot be empty")
	}

	team := futchain.GetRootAsTeam(data, 0)

	return &datasource.Team{
		ID:       int(team.Id()),
		Score:    int(team.Score()),
		Name:     string(team.Name()),
		LongName: string(team.LongName()),
	}, nil
}

// EncodeStatus encodes a Status struct to FlatBuffers binary format
func EncodeStatus(status *datasource.Status) ([]byte, error) {
	if status == nil {
		return nil, errors.New("status cannot be nil")
	}

	builder := flatbuffers.NewBuilder(1024)

	// Create the Status table
	futchain.StatusStart(builder)
	futchain.StatusAddUtcTime(builder, status.UtcTime.Unix())
	futchain.StatusAddPeriodLength(builder, int32(status.PeriodLength))
	futchain.StatusAddStarted(builder, status.Started)
	futchain.StatusAddCancelled(builder, status.Cancelled)
	futchain.StatusAddFinished(builder, status.Finished)
	statusOffset := futchain.StatusEnd(builder)

	builder.Finish(statusOffset)
	return builder.FinishedBytes(), nil
}

// DecodeStatus decodes FlatBuffers binary data to a Status struct
func DecodeStatus(data []byte) (*datasource.Status, error) {
	if len(data) == 0 {
		return nil, errors.New("data cannot be empty")
	}

	status := futchain.GetRootAsStatus(data, 0)

	return &datasource.Status{
		UtcTime:      time.Unix(status.UtcTime(), 0),
		PeriodLength: int(status.PeriodLength()),
		Started:      status.Started(),
		Cancelled:    status.Cancelled(),
		Finished:     status.Finished(),
	}, nil
}

// EncodeMatch encodes a Match struct to FlatBuffers binary format
func EncodeMatch(match *datasource.Match) ([]byte, error) {
	if match == nil {
		return nil, errors.New("match cannot be nil")
	}

	builder := flatbuffers.NewBuilder(2048)

	// Create string offsets
	timeOffset := builder.CreateString(match.Time)
	tournamentStageOffset := builder.CreateString(match.TournamentStage)

	// Encode nested Team structs directly
	homeNameOffset := builder.CreateString(match.Home.Name)
	homeLongNameOffset := builder.CreateString(match.Home.LongName)
	futchain.TeamStart(builder)
	futchain.TeamAddId(builder, int32(match.Home.ID))
	futchain.TeamAddScore(builder, int32(match.Home.Score))
	futchain.TeamAddName(builder, homeNameOffset)
	futchain.TeamAddLongName(builder, homeLongNameOffset)
	homeOffset := futchain.TeamEnd(builder)

	awayNameOffset := builder.CreateString(match.Away.Name)
	awayLongNameOffset := builder.CreateString(match.Away.LongName)
	futchain.TeamStart(builder)
	futchain.TeamAddId(builder, int32(match.Away.ID))
	futchain.TeamAddScore(builder, int32(match.Away.Score))
	futchain.TeamAddName(builder, awayNameOffset)
	futchain.TeamAddLongName(builder, awayLongNameOffset)
	awayOffset := futchain.TeamEnd(builder)

	// Encode Status struct directly
	futchain.StatusStart(builder)
	futchain.StatusAddUtcTime(builder, match.Status.UtcTime.Unix())
	futchain.StatusAddPeriodLength(builder, int32(match.Status.PeriodLength))
	futchain.StatusAddStarted(builder, match.Status.Started)
	futchain.StatusAddCancelled(builder, match.Status.Cancelled)
	futchain.StatusAddFinished(builder, match.Status.Finished)
	statusOffset := futchain.StatusEnd(builder)

	// Handle eliminated team ID (convert any to int32)
	eliminatedTeamID := int32(-1) // Default to -1 for null
	if match.EliminatedTeamID != nil {
		if id, ok := match.EliminatedTeamID.(int); ok {
			eliminatedTeamID = int32(id)
		} else if id, ok := match.EliminatedTeamID.(int32); ok {
			eliminatedTeamID = id
		}
	}

	// Create the Match table
	futchain.MatchStart(builder)
	futchain.MatchAddId(builder, int32(match.ID))
	futchain.MatchAddLeagueId(builder, int32(match.LeagueID))
	futchain.MatchAddTime(builder, timeOffset)
	futchain.MatchAddHome(builder, homeOffset)
	futchain.MatchAddAway(builder, awayOffset)
	futchain.MatchAddEliminatedTeamId(builder, eliminatedTeamID)
	futchain.MatchAddStatusId(builder, int32(match.StatusID))
	futchain.MatchAddTournamentStage(builder, tournamentStageOffset)
	futchain.MatchAddStatus(builder, statusOffset)
	futchain.MatchAddTimeTs(builder, match.TimeTS)
	matchOffset := futchain.MatchEnd(builder)

	builder.Finish(matchOffset)
	return builder.FinishedBytes(), nil
}

// DecodeMatch decodes FlatBuffers binary data to a Match struct
func DecodeMatch(data []byte) (*datasource.Match, error) {
	if len(data) == 0 {
		return nil, errors.New("data cannot be empty")
	}

	match := futchain.GetRootAsMatch(data, 0)

	// Decode nested structs using the correct FlatBuffers API
	var home datasource.Team
	var away datasource.Team
	var status datasource.Status

	// Get nested objects
	homeObj := match.Home(nil)
	if homeObj != nil {
		home = datasource.Team{
			ID:       int(homeObj.Id()),
			Score:    int(homeObj.Score()),
			Name:     string(homeObj.Name()),
			LongName: string(homeObj.LongName()),
		}
	}

	awayObj := match.Away(nil)
	if awayObj != nil {
		away = datasource.Team{
			ID:       int(awayObj.Id()),
			Score:    int(awayObj.Score()),
			Name:     string(awayObj.Name()),
			LongName: string(awayObj.LongName()),
		}
	}

	statusObj := match.Status(nil)
	if statusObj != nil {
		status = datasource.Status{
			UtcTime:      time.Unix(statusObj.UtcTime(), 0),
			PeriodLength: int(statusObj.PeriodLength()),
			Started:      statusObj.Started(),
			Cancelled:    statusObj.Cancelled(),
			Finished:     statusObj.Finished(),
		}
	}

	// Handle eliminated team ID
	var eliminatedTeamID any
	if match.EliminatedTeamId() != -1 {
		eliminatedTeamID = int(match.EliminatedTeamId())
	}

	return &datasource.Match{
		ID:               int(match.Id()),
		LeagueID:         int(match.LeagueId()),
		Time:             string(match.Time()),
		Home:             home,
		Away:             away,
		EliminatedTeamID: eliminatedTeamID,
		StatusID:         int(match.StatusId()),
		TournamentStage:  string(match.TournamentStage()),
		Status:           status,
		TimeTS:           match.TimeTs(),
	}, nil
}

// EncodeLeague encodes a League struct to FlatBuffers binary format
func EncodeLeague(league *datasource.League) ([]byte, error) {
	if league == nil {
		return nil, errors.New("league cannot be nil")
	}

	builder := flatbuffers.NewBuilder(1024)

	// Create string offsets
	groupNameOffset := builder.CreateString(league.GroupName)
	ccodeOffset := builder.CreateString(league.Ccode)
	nameOffset := builder.CreateString(league.Name)

	// Create the League table
	futchain.LeagueStart(builder)
	futchain.LeagueAddIsGroup(builder, league.IsGroup)
	futchain.LeagueAddGroupName(builder, groupNameOffset)
	futchain.LeagueAddCcode(builder, ccodeOffset)
	futchain.LeagueAddId(builder, int32(league.ID))
	futchain.LeagueAddPrimaryId(builder, int32(league.PrimaryID))
	futchain.LeagueAddName(builder, nameOffset)
	leagueOffset := futchain.LeagueEnd(builder)

	builder.Finish(leagueOffset)
	return builder.FinishedBytes(), nil
}

// DecodeLeague decodes FlatBuffers binary data to a League struct
func DecodeLeague(data []byte) (*datasource.League, error) {
	if len(data) == 0 {
		return nil, errors.New("data cannot be empty")
	}

	league := futchain.GetRootAsLeague(data, 0)

	return &datasource.League{
		IsGroup:   league.IsGroup(),
		GroupName: string(league.GroupName()),
		Ccode:     string(league.Ccode()),
		ID:        int(league.Id()),
		PrimaryID: int(league.PrimaryId()),
		Name:      string(league.Name()),
		Matches:   []datasource.Match{}, // Empty matches slice
	}, nil
}

// GetBinarySize returns the size of the binary data
func GetBinarySize(data []byte) int {
	return len(data)
}
