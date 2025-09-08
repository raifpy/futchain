package futchain

import (
	"fmt"
	"math/big"

	_ "embed"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/raifpy/futchain/x/futchain/keeper"
	futchaintypes "github.com/raifpy/futchain/x/futchain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/evm/x/vm/statedb"
)

var (
	_ vm.PrecompiledContract = (*FutchainEvmBridge)(nil)
)

type FutchainEvmBridge struct {
	keeper *keeper.Keeper
	abi    abi.ABI
}

func NewFutchainEvmBridge(keeper *keeper.Keeper) (*FutchainEvmBridge, error) {
	return &FutchainEvmBridge{keeper: keeper, abi: keeper.ABI}, nil
}

func (f *FutchainEvmBridge) Address() common.Address {
	return common.HexToAddress(futchaintypes.FutchainPrecompileAddress)
}

func (f *FutchainEvmBridge) RequiredGas(_ []byte) uint64 {
	//return 10000 // Set a reasonable gas cost for your precompile
	return 0
}

func (f *FutchainEvmBridge) Run(evm *vm.EVM, contract *vm.Contract, ronly bool) ([]byte, error) {
	if len(contract.Input) < 4 {
		return nil, vm.ErrExecutionReverted
	}

	method, err := f.abi.MethodById(contract.Input[:4])
	if err != nil {
		fmt.Printf("Method lookup error: %v\n", err)
		return nil, err
	}

	// Parse input arguments
	argsBz := contract.Input[4:]
	args, err := method.Inputs.Unpack(argsBz)
	if err != nil {
		fmt.Printf("[evm] Input unpacking error: %v\n", err)
		return nil, err
	}

	// Get the proper SDK context from StateDB
	stateDB, ok := evm.StateDB.(*statedb.StateDB)
	if !ok {
		return nil, fmt.Errorf("invalid StateDB type")
	}

	ctx := stateDB.GetContext()

	switch method.Name {
	case "getMatch":
		return f.handleGetMatch(ctx, method, args)
	case "getLeague":
		return f.handleGetLeague(ctx, method, args)
	case "getTeam":
		return f.handleGetTeam(ctx, method, args)
	case "getUnfinishedMatches":
		return f.handleGetUnfinishedMatches(ctx, method, args)
	}

	return nil, fmt.Errorf("method %s not implemented", method.Name)
}

// handleGetMatch handles the getMatch function call
func (f *FutchainEvmBridge) handleGetMatch(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments for getMatch")
	}

	matchId, ok := args[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid matchId type")
	}

	match, err := f.keeper.GetMatch(ctx, int(matchId.Int64()))
	if err != nil {
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	// Create the struct tuple for MatchData
	matchData := struct {
		Id        *big.Int
		LeagueId  *big.Int
		Name      string
		Time      string
		HomeId    *big.Int
		AwayId    *big.Int
		HomeScore *big.Int
		AwayScore *big.Int
		HomeName  string
		AwayName  string
		Started   bool
		Finished  bool
		Cancelled bool
	}{
		Id:        big.NewInt(int64(match.ID)),
		LeagueId:  big.NewInt(int64(match.LeagueID)),
		Name:      match.Home.Name + " - " + match.Away.Name,
		Time:      match.Time, // example time:  "09.09.2025 20:45"
		HomeId:    big.NewInt(int64(match.Home.ID)),
		AwayId:    big.NewInt(int64(match.Away.ID)),
		HomeScore: big.NewInt(int64(match.Home.Score)),
		AwayScore: big.NewInt(int64(match.Away.Score)),
		HomeName:  match.Home.Name,
		AwayName:  match.Away.Name,
		Started:   match.Status.Started,
		Finished:  match.Status.Finished,
		Cancelled: match.Status.Cancelled,
	}

	return method.Outputs.Pack(matchData)
}

// handleGetLeague handles the getLeague function call
func (f *FutchainEvmBridge) handleGetLeague(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments for getLeague")
	}

	leagueId, ok := args[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid leagueId type")
	}

	league, err := f.keeper.GetLeague(ctx, int(leagueId.Int64()))
	if err != nil {
		return nil, fmt.Errorf("failed to get league: %w", err)
	}

	// Create the struct tuple for LeagueData
	leagueData := struct {
		Id        *big.Int
		Name      string
		GroupName string
	}{
		Id:        big.NewInt(int64(league.ID)),
		Name:      league.Name,
		GroupName: league.GroupName,
	}

	return method.Outputs.Pack(leagueData)
}

// handleGetTeam handles the getTeam function call
func (f *FutchainEvmBridge) handleGetTeam(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments for getTeam")
	}

	teamId, ok := args[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("invalid teamId type")
	}

	team, err := f.keeper.GetTeam(ctx, int(teamId.Int64()))
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	// Create the struct tuple for TeamData
	teamData := struct {
		Id   *big.Int
		Name string
	}{
		Id:   big.NewInt(int64(team.ID)),
		Name: team.Name,
	}

	return method.Outputs.Pack(teamData)
}

// handleGetUnfinishedMatches handles the getUnfinishedMatches function call
func (f *FutchainEvmBridge) handleGetUnfinishedMatches(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("invalid number of arguments for getUnfinishedMatches")
	}

	ids, err := f.keeper.ListUnfinishedMatches(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get unfinished matches: %w", err)
	}

	// Convert to []*big.Int for ABI encoding
	bigIntIds := make([]*big.Int, len(ids))
	for i, id := range ids {
		bigIntIds[i] = big.NewInt(int64(id))
	}

	return method.Outputs.Pack(bigIntIds)
}
