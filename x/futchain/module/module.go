package futchain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"github.com/raifpy/futchainevm/x/futchain/keeper"
	"github.com/raifpy/futchainevm/x/futchain/keeper/datasource"
	"github.com/raifpy/futchainevm/x/futchain/types"
)

var (
	_ module.AppModuleBasic = (*AppModule)(nil)
	_ module.AppModule      = (*AppModule)(nil)
	_ module.HasGenesis     = (*AppModule)(nil)

	_ appmodule.AppModule       = (*AppModule)(nil)
	_ appmodule.HasBeginBlocker = (*AppModule)(nil)
	_ appmodule.HasEndBlocker   = (*AppModule)(nil)
)

// AppModule implements the AppModule interface that defines the inter-dependent methods that modules need to implement
type AppModule struct {
	cdc        codec.Codec
	keeper     keeper.Keeper
	authKeeper types.AuthKeeper
	bankKeeper types.BankKeeper
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	authKeeper types.AuthKeeper,
	bankKeeper types.BankKeeper,
) AppModule {
	return AppModule{
		cdc:        cdc,
		keeper:     keeper,
		authKeeper: authKeeper,
		bankKeeper: bankKeeper,
	}
}

// IsAppModule implements the appmodule.AppModule interface.
func (AppModule) IsAppModule() {}

// Name returns the name of the module as a string.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the amino codec
func (AppModule) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModule) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(clientCtx.CmdContext, mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// RegisterInterfaces registers a module's interface types and their concrete implementations as proto.Message.
func (AppModule) RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registrar)
}

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC queries
func (am AppModule) RegisterServices(registrar grpc.ServiceRegistrar) error {
	types.RegisterMsgServer(registrar, keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(registrar, keeper.NewQueryServerImpl(am.keeper))

	return nil
}

// DefaultGenesis returns a default GenesisState for the module, marshalled to json.RawMessage.
// The default GenesisState need to be defined by the module developer and is primarily used for testing.
func (am AppModule) DefaultGenesis(codec.JSONCodec) json.RawMessage {
	return am.cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis used to validate the GenesisState, given in its json.RawMessage form.
func (am AppModule) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := am.cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return genState.Validate()
}

// InitGenesis performs the module's genesis initialization. It returns no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, gs json.RawMessage) {
	var genState types.GenesisState
	// Initialize global index to index in genesis state
	if err := am.cdc.UnmarshalJSON(gs, &genState); err != nil {
		panic(fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err))
	}

	if err := am.keeper.InitGenesis(ctx, genState); err != nil {
		panic(fmt.Errorf("failed to initialize %s genesis state: %w", types.ModuleName, err))
	}
}

// ExportGenesis returns the module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONCodec) json.RawMessage {
	genState, err := am.keeper.ExportGenesis(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to export %s genesis state: %w", types.ModuleName, err))
	}

	bz, err := am.cdc.MarshalJSON(genState)
	if err != nil {
		panic(fmt.Errorf("failed to marshal %s genesis state: %w", types.ModuleName, err))
	}

	return bz
}

// ConsensusVersion is a sequence number for state-breaking change of the module.
// It should be incremented on each consensus-breaking change introduced by the module.
// To avoid wrong/empty versions, the initial version should be set to 1.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock contains the logic that is automatically triggered at the beginning of each block.
// The begin block implementation is optional.
func (am AppModule) BeginBlock(goCtx context.Context) error {

	log.Printf("\n\nRunning BeginBlock\n\n")
	ctx := sdk.UnwrapSDKContext(goCtx)

	params, err := am.keeper.Params.Get(goCtx)
	if err != nil {
		ctx.Logger().Error("error getting genesis params", "error", err)
		return err
	}

	if ctx.BlockHeight()%params.FetchModulo != 0 {
		ctx.Logger().Debug("block height is not divisible by the fetch modulo. skipping fetch", "block height", ctx.BlockHeight(), "fetch modulo", params.FetchModulo)
		return nil // we don't fetch data if the block height is not divisible by the fetch modulo
	}

	ctx.Logger().Info("fetching data", "block height", ctx.BlockHeight(), "fetch modulo", params.FetchModulo)
	result, err := am.keeper.Datasource.Fetch(goCtx, datasource.WithLogger(ctx.Logger().With("source", "datasource")), datasource.WithTimezone(params.Timezone))
	if err != nil {
		ctx.Logger().Error("failed to fetch data", "error", err)
		return err
	}

	for _, l := range result {

		saved, err := am.keeper.SaveLeagueIfNotExists(goCtx, l)
		if err != nil {
			ctx.Logger().Error("failed to save league to the store", "error", err)
			continue
		}

		if saved {
			// event that we have detected a new league
			ctx.Logger().Info("detected a new league", "league", l.Name, "id", l.ID, "group", l.GroupName, "event", "new_league")
			ctx.EventManager().EmitEvent(sdk.NewEvent("new_league", sdk.NewAttribute("league", l.Name), sdk.NewAttribute("id", strconv.Itoa(l.ID))))
			//TODO: set TypedEvent
		}
		for _, m := range l.Matches {
			// save teams if not exists
			_, err := am.keeper.SaveTeamIfNotExists(goCtx, m.Home)
			if err != nil {
				ctx.Logger().Error("failed to save home team to the store", "error", err, "team", m.Home.Name, "id", m.Home.ID)
			}
			_, err = am.keeper.SaveTeamIfNotExists(goCtx, m.Away)
			if err != nil {
				ctx.Logger().Error("failed to save away team to the store", "error", err, "team", m.Away.Name, "id", m.Away.ID)
			}

			saved, err := am.keeper.SaveMatchIfNotExists(goCtx, m)
			if err != nil {
				ctx.Logger().Error("failed to save match to the store", "error", err, "match", m.ID, "league_id", m.LeagueID, "home_id", m.Home.ID, "away_id", m.Away.ID)
				continue
			}

			if saved {
				ctx.Logger().Info("detected a new match", "match", m.ID, "event", "new_match")
				ctx.EventManager().EmitEvent(sdk.NewEvent("new_match", sdk.NewAttribute("id", strconv.Itoa(m.ID)), sdk.NewAttribute("league_id", strconv.Itoa(m.LeagueID)), sdk.NewAttribute("match", m.Home.Name+"/"+m.Away.Name), sdk.NewAttribute("home_id", strconv.Itoa(m.Home.ID)), sdk.NewAttribute("away_id", strconv.Itoa(m.Away.ID)), sdk.NewAttribute("event", "new_match")))
			} else {
				//compare for match updates
				oldmatch, err := am.keeper.GetMatch(goCtx, m.ID)
				if err != nil {
					// unexpected error
					ctx.Logger().Error("failed to get match from the store", "error", err, "when", "compare_match_updates", "match", m.ID, "league_id", m.LeagueID, "home_id", m.Home.ID, "away_id", m.Away.ID)
					continue
				}

				if pri := m.Compare(oldmatch); pri != datasource.PriorityNoChanges {

					err := am.keeper.SetMatch(goCtx, m)
					if err != nil {
						ctx.Logger().Error("failed to set match to the store", "error", err, "match", m.ID, "league_id", m.LeagueID, "home_id", m.Home.ID, "away_id", m.Away.ID)
						continue
					}

					// match has changed.
					if pri >= datasource.MinimumEventPriority {

						// we will emit an event pri.EventName()
						ctx.Logger().Info("match has changed", "match", m.ID, "event", pri.EventName())
						ctx.EventManager().EmitEvent(sdk.NewEvent(pri.EventName(), sdk.NewAttribute("id", strconv.Itoa(m.ID)), sdk.NewAttribute("league_id", strconv.Itoa(m.LeagueID)), sdk.NewAttribute("match", m.Home.Name+"/"+m.Away.Name), sdk.NewAttribute("home_id", strconv.Itoa(m.Home.ID)), sdk.NewAttribute("away_id", strconv.Itoa(m.Away.ID)), sdk.NewAttribute("event", pri.EventName())))
					}
				} else {
					ctx.Logger().Debug("match has no changes", "match", m.ID, "league_id", m.LeagueID, "home_id", m.Home.ID, "away_id", m.Away.ID)
				}

			}
		}

	}

	return nil
}

// EndBlock contains the logic that is automatically triggered at the end of each block.
// The end block implementation is optional.
func (am AppModule) EndBlock(_ context.Context) error {
	return nil
}
