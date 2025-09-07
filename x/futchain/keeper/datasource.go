package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/raifpy/futchainevm/x/futchain/keeper/datasource"
	"github.com/raifpy/futchainevm/x/futchain/keeper/datasource/flatbuffers"
)

type DatasourceConfig struct {
	ApiURL  string
	Headers map[string]string
}

func (k *Keeper) TeamKey(id int) []byte {

	return append(TeamKey, sdk.Uint64ToBigEndian(uint64(id))...)
}

func (k *Keeper) MatchKey(id int) []byte {
	return append(MatchKey, sdk.Uint64ToBigEndian(uint64(id))...)
}

func (k *Keeper) LeagueKey(id int) []byte {
	return append(LeagueKey, sdk.Uint64ToBigEndian(uint64(id))...)
}

func (k *Keeper) SaveTeamIfNotExists(ctx context.Context, team datasource.Team) (bool, error) {
	key := k.TeamKey(team.ID)
	if ok, err := k.storeService.OpenKVStore(ctx).Has(key); ok {
		return false, nil
	} else if err != nil {
		return false, err
	}

	buf, err := flatbuffers.NewTeamEncoder().EncodeToBinary(&team)
	if err != nil {
		return false, err
	}
	//TODO: use proto marshaler
	return true, k.storeService.OpenKVStore(ctx).Set(key, buf)
}

func (k *Keeper) SaveMatchIfNotExists(ctx context.Context, match datasource.Match) (bool, error) {
	key := k.MatchKey(match.ID)
	if ok, err := k.storeService.OpenKVStore(ctx).Has(key); ok {
		return false, nil
	} else if err != nil {
		return false, err
	}

	buf, err := flatbuffers.NewMatchEncoder().EncodeToBinary(&match)
	if err != nil {
		return false, err
	}
	//TODO: use proto marshaler
	return true, k.storeService.OpenKVStore(ctx).Set(key, buf)
}

func (k *Keeper) SaveLeagueIfNotExists(ctx context.Context, league datasource.League) (bool, error) {
	key := k.LeagueKey(league.ID)
	if ok, err := k.storeService.OpenKVStore(ctx).Has(key); ok {
		return false, nil
	} else if err != nil {
		return false, err
	}

	buf, err := flatbuffers.NewLeagueEncoder().EncodeToBinary(&league)
	if err != nil {
		return false, err
	}
	//TODO: use proto marshaler
	return true, k.storeService.OpenKVStore(ctx).Set(key, buf)
}

func (k *Keeper) GetMatch(ctx context.Context, id int) (*datasource.Match, error) {
	key := k.MatchKey(id)
	buf, err := k.storeService.OpenKVStore(ctx).Get(key)
	if err != nil {
		return nil, err
	}
	match, err := flatbuffers.NewMatchEncoder().DecodeFromBinary(buf)
	if err != nil {
		return nil, err
	}

	home, err := k.GetTeam(ctx, match.Home.ID)
	if err != nil {
		return nil, err
	}
	home.Score = match.Home.Score
	home.ID = match.Home.ID
	match.Home = *home

	away, err := k.GetTeam(ctx, match.Away.ID)
	if err != nil {
		return nil, err
	}
	away.Score = match.Away.Score
	away.ID = match.Away.ID
	match.Away = *away
	return match, nil

}

func (k *Keeper) SetMatch(ctx context.Context, match datasource.Match) error {
	key := k.MatchKey(match.ID)
	buf, err := flatbuffers.NewMatchEncoder().EncodeToBinary(&match)
	if err != nil {
		return err
	}
	return k.storeService.OpenKVStore(ctx).Set(key, buf)
}

func (k *Keeper) GetTeam(ctx context.Context, id int) (*datasource.Team, error) {
	key := k.TeamKey(id)
	buf, err := k.storeService.OpenKVStore(ctx).Get(key)
	if err != nil {
		return nil, err
	}
	return flatbuffers.NewTeamEncoder().DecodeFromBinary(buf)
}
