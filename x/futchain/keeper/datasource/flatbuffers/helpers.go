package flatbuffers

import (
	"encoding/hex"
	"fmt"

	"github.com/raifpy/futchain/x/futchain/keeper/datasource"
)

// TeamEncoder provides blockchain-style encoding for Team structs
type TeamEncoder struct{}

// NewTeamEncoder creates a new TeamEncoder instance
func NewTeamEncoder() *TeamEncoder {
	return &TeamEncoder{}
}

// EncodeToBinary encodes a Team to binary format for blockchain storage
func (te *TeamEncoder) EncodeToBinary(team *datasource.Team) ([]byte, error) {
	return EncodeTeam(team)
}

// DecodeFromBinary decodes binary data back to a Team struct
func (te *TeamEncoder) DecodeFromBinary(data []byte) (*datasource.Team, error) {
	return DecodeTeam(data)
}

// EncodeToHex encodes a Team to hex string (useful for debugging and logging)
func (te *TeamEncoder) EncodeToHex(team *datasource.Team) (string, error) {
	data, err := EncodeTeam(team)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

// DecodeFromHex decodes hex string back to a Team struct
func (te *TeamEncoder) DecodeFromHex(hexStr string) (*datasource.Team, error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %w", err)
	}
	return DecodeTeam(data)
}

// GetSize returns the binary size of a Team when encoded
func (te *TeamEncoder) GetSize(team *datasource.Team) (int, error) {
	data, err := EncodeTeam(team)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

// MatchEncoder provides blockchain-style encoding for Match structs
type MatchEncoder struct{}

// NewMatchEncoder creates a new MatchEncoder instance
func NewMatchEncoder() *MatchEncoder {
	return &MatchEncoder{}
}

// EncodeToBinary encodes a Match to binary format for blockchain storage
func (me *MatchEncoder) EncodeToBinary(match *datasource.Match) ([]byte, error) {
	return EncodeMatch(match)
}

// DecodeFromBinary decodes binary data back to a Match struct
func (me *MatchEncoder) DecodeFromBinary(data []byte) (*datasource.Match, error) {
	return DecodeMatch(data)
}

// EncodeToHex encodes a Match to hex string (useful for debugging and logging)
func (me *MatchEncoder) EncodeToHex(match *datasource.Match) (string, error) {
	data, err := EncodeMatch(match)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

// DecodeFromHex decodes hex string back to a Match struct
func (me *MatchEncoder) DecodeFromHex(hexStr string) (*datasource.Match, error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %w", err)
	}
	return DecodeMatch(data)
}

// GetSize returns the binary size of a Match when encoded
func (me *MatchEncoder) GetSize(match *datasource.Match) (int, error) {
	data, err := EncodeMatch(match)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

// LeagueEncoder provides blockchain-style encoding for League structs
type LeagueEncoder struct{}

// NewLeagueEncoder creates a new LeagueEncoder instance
func NewLeagueEncoder() *LeagueEncoder {
	return &LeagueEncoder{}
}

// EncodeToBinary encodes a League to binary format for blockchain storage
func (le *LeagueEncoder) EncodeToBinary(league *datasource.League) ([]byte, error) {
	return EncodeLeague(league)
}

// DecodeFromBinary decodes binary data back to a League struct
func (le *LeagueEncoder) DecodeFromBinary(data []byte) (*datasource.League, error) {
	return DecodeLeague(data)
}

// EncodeToHex encodes a League to hex string (useful for debugging and logging)
func (le *LeagueEncoder) EncodeToHex(league *datasource.League) (string, error) {
	data, err := EncodeLeague(league)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

// DecodeFromHex decodes hex string back to a League struct
func (le *LeagueEncoder) DecodeFromHex(hexStr string) (*datasource.League, error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %w", err)
	}
	return DecodeLeague(data)
}

// GetSize returns the binary size of a League when encoded
func (le *LeagueEncoder) GetSize(league *datasource.League) (int, error) {
	data, err := EncodeLeague(league)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

// UniversalEncoder provides encoding for all struct types
type UniversalEncoder struct {
	teamEncoder   *TeamEncoder
	matchEncoder  *MatchEncoder
	leagueEncoder *LeagueEncoder
}

// NewUniversalEncoder creates a new UniversalEncoder instance
func NewUniversalEncoder() *UniversalEncoder {
	return &UniversalEncoder{
		teamEncoder:   NewTeamEncoder(),
		matchEncoder:  NewMatchEncoder(),
		leagueEncoder: NewLeagueEncoder(),
	}
}

// EncodeTeam encodes a Team using the universal encoder
func (ue *UniversalEncoder) EncodeTeam(team *datasource.Team) ([]byte, error) {
	return ue.teamEncoder.EncodeToBinary(team)
}

// DecodeTeam decodes a Team using the universal encoder
func (ue *UniversalEncoder) DecodeTeam(data []byte) (*datasource.Team, error) {
	return ue.teamEncoder.DecodeFromBinary(data)
}

// EncodeMatch encodes a Match using the universal encoder
func (ue *UniversalEncoder) EncodeMatch(match *datasource.Match) ([]byte, error) {
	return ue.matchEncoder.EncodeToBinary(match)
}

// DecodeMatch decodes a Match using the universal encoder
func (ue *UniversalEncoder) DecodeMatch(data []byte) (*datasource.Match, error) {
	return ue.matchEncoder.DecodeFromBinary(data)
}

// EncodeLeague encodes a League using the universal encoder
func (ue *UniversalEncoder) EncodeLeague(league *datasource.League) ([]byte, error) {
	return ue.leagueEncoder.EncodeToBinary(league)
}

// DecodeLeague decodes a League using the universal encoder
func (ue *UniversalEncoder) DecodeLeague(data []byte) (*datasource.League, error) {
	return ue.leagueEncoder.DecodeFromBinary(data)
}
