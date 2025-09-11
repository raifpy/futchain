pragma solidity >=0.8.17;

address constant FUTAPP_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000807;

// Futchain Data Structures
struct MatchData {
    uint256 id;
    uint256 leagueId;
    string time;
    string minute;
    uint256 homeId;
    uint256 awayId;
    uint256 homeScore;
    uint256 awayScore;
    string homeName;
    string awayName;
    bool started;
    bool finished;
    bool cancelled;
}

struct LeagueData {
    uint256 id;
    string name;
    string groupName;
}

struct TeamData {
    uint256 id;
    string name;
}

// Futchain Interface Contract
interface FutI {
    /// @notice Get match details by ID
    /// @param matchId The match ID to query
    /// @return match The match data structure
    function getMatch(uint256 matchId) external view returns (MatchData memory);
    
    /// @notice Get league details by ID
    /// @param leagueId The league ID to query
    /// @return league The league data structure
    function getLeague(uint256 leagueId) external view returns (LeagueData memory);
    
    /// @notice Get team details by ID
    /// @param teamId The team ID to query
    /// @return team The team data structure
    function getTeam(uint256 teamId) external view returns (TeamData memory);
    
    /// @notice Get list of unfinished match IDs
    /// @return matchIds Array of unfinished match IDs
    function getUnfinishedMatches() external view returns (uint256[] memory);
}

// Futchain Precompile Instance
FutI constant FUTCHAIN = FutI(FUTAPP_PRECOMPILE_ADDRESS);