package flatbuffers

import (
	"testing"
	"time"

	"github.com/raifpy/futchain/x/futchain/keeper/datasource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTeamEncoding(t *testing.T) {
	// Create a test team
	team := &datasource.Team{
		ID:       123,
		Score:    2,
		Name:     "Arsenal",
		LongName: "Arsenal Football Club",
	}

	// Test encoding
	encoder := NewTeamEncoder()
	data, err := encoder.EncodeToBinary(team)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Test decoding
	decodedTeam, err := encoder.DecodeFromBinary(data)
	require.NoError(t, err)
	assert.Equal(t, team.ID, decodedTeam.ID)
	assert.Equal(t, team.Score, decodedTeam.Score)
	assert.Equal(t, team.Name, decodedTeam.Name)
	assert.Equal(t, team.LongName, decodedTeam.LongName)

	// Test hex encoding
	hexStr, err := encoder.EncodeToHex(team)
	require.NoError(t, err)
	assert.NotEmpty(t, hexStr)

	// Test hex decoding
	decodedFromHex, err := encoder.DecodeFromHex(hexStr)
	require.NoError(t, err)
	assert.Equal(t, team, decodedFromHex)

	// Test size calculation
	size, err := encoder.GetSize(team)
	require.NoError(t, err)
	assert.Equal(t, len(data), size)
}

func TestMatchEncoding(t *testing.T) {
	// Create test teams
	homeTeam := &datasource.Team{
		ID:       1,
		Score:    2,
		Name:     "Arsenal",
		LongName: "Arsenal Football Club",
	}

	awayTeam := &datasource.Team{
		ID:       2,
		Score:    1,
		Name:     "Chelsea",
		LongName: "Chelsea Football Club",
	}

	// Create test status
	status := &datasource.Status{
		UtcTime:      time.Now(),
		PeriodLength: 90,
		Started:      true,
		Cancelled:    false,
		Finished:     false,
	}

	// Create test match
	match := &datasource.Match{
		ID:               1001,
		LeagueID:         1,
		Time:             "15:00",
		Home:             *homeTeam,
		Away:             *awayTeam,
		EliminatedTeamID: nil,
		StatusID:         1,
		TournamentStage:  "Regular Season",
		Status:           *status,
		TimeTS:           time.Now().Unix(),
	}

	// Test encoding
	encoder := NewMatchEncoder()
	data, err := encoder.EncodeToBinary(match)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Test decoding
	decodedMatch, err := encoder.DecodeFromBinary(data)
	require.NoError(t, err)
	assert.Equal(t, match.ID, decodedMatch.ID)
	assert.Equal(t, match.LeagueID, decodedMatch.LeagueID)
	assert.Equal(t, match.Time, decodedMatch.Time)
	assert.Equal(t, match.Home.ID, decodedMatch.Home.ID) // Only ID is preserved
	assert.Equal(t, match.Away.ID, decodedMatch.Away.ID) // Only ID is preserved
	assert.Equal(t, match.StatusID, decodedMatch.StatusID)
	assert.Equal(t, match.TournamentStage, decodedMatch.TournamentStage)
	assert.Equal(t, match.TimeTS, decodedMatch.TimeTS)

	// Test hex encoding
	hexStr, err := encoder.EncodeToHex(match)
	require.NoError(t, err)
	assert.NotEmpty(t, hexStr)

	// Test size calculation
	size, err := encoder.GetSize(match)
	require.NoError(t, err)
	assert.Equal(t, len(data), size)
}

func TestLeagueEncoding(t *testing.T) {
	// Create test teams
	homeTeam := &datasource.Team{
		ID:       1,
		Score:    2,
		Name:     "Arsenal",
		LongName: "Arsenal Football Club",
	}

	awayTeam := &datasource.Team{
		ID:       2,
		Score:    1,
		Name:     "Chelsea",
		LongName: "Chelsea Football Club",
	}

	// Create test status
	status := &datasource.Status{
		UtcTime:      time.Now(),
		PeriodLength: 90,
		Started:      true,
		Cancelled:    false,
		Finished:     false,
	}

	// Create test match
	match := &datasource.Match{
		ID:               1001,
		LeagueID:         1,
		Time:             "15:00",
		Home:             *homeTeam,
		Away:             *awayTeam,
		EliminatedTeamID: nil,
		StatusID:         1,
		TournamentStage:  "Regular Season",
		Status:           *status,
		TimeTS:           time.Now().Unix(),
	}

	// Create test league
	league := &datasource.League{
		IsGroup:   false,
		GroupName: "",
		Ccode:     "GBR",
		ID:        1,
		PrimaryID: 1,
		Name:      "Premier League",
		Matches:   []datasource.Match{*match},
	}

	// Test encoding
	encoder := NewLeagueEncoder()
	data, err := encoder.EncodeToBinary(league)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Test decoding
	decodedLeague, err := encoder.DecodeFromBinary(data)
	require.NoError(t, err)
	assert.Equal(t, league.IsGroup, decodedLeague.IsGroup)
	assert.Equal(t, league.GroupName, decodedLeague.GroupName)
	assert.Equal(t, league.Ccode, decodedLeague.Ccode)
	assert.Equal(t, league.ID, decodedLeague.ID)
	assert.Equal(t, league.PrimaryID, decodedLeague.PrimaryID)
	assert.Equal(t, league.Name, decodedLeague.Name)
	assert.Len(t, decodedLeague.Matches, 0) // No matches should be encoded/decoded

	// Test hex encoding
	hexStr, err := encoder.EncodeToHex(league)
	require.NoError(t, err)
	assert.NotEmpty(t, hexStr)

	// Test size calculation
	size, err := encoder.GetSize(league)
	require.NoError(t, err)
	assert.Equal(t, len(data), size)
}

func TestUniversalEncoder(t *testing.T) {
	// Create test data
	team := &datasource.Team{
		ID:       123,
		Score:    2,
		Name:     "Arsenal",
		LongName: "Arsenal Football Club",
	}

	// Test universal encoder
	encoder := NewUniversalEncoder()

	// Test team encoding
	data, err := encoder.EncodeTeam(team)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	decodedTeam, err := encoder.DecodeTeam(data)
	require.NoError(t, err)
	assert.Equal(t, team, decodedTeam)
}

func TestErrorHandling(t *testing.T) {
	encoder := NewTeamEncoder()

	// Test nil input
	_, err := encoder.EncodeToBinary(nil)
	assert.Error(t, err)

	// Test empty data
	_, err = encoder.DecodeFromBinary([]byte{})
	assert.Error(t, err)

	// Test invalid hex
	_, err = encoder.DecodeFromHex("invalid hex")
	assert.Error(t, err)
}
