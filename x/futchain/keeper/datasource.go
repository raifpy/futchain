package keeper

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/raifpy/futchain/x/futchain/keeper/datasource"
	"github.com/raifpy/futchain/x/futchain/keeper/datasource/flatbuffers"
)

type DatasourceConfig struct {
	ApiURL  string
	Headers map[string]string
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

func (k *Keeper) GetLeague(ctx context.Context, id int) (*datasource.League, error) {
	key := k.LeagueKey(id)
	buf, err := k.storeService.OpenKVStore(ctx).Get(key)
	if err != nil {
		return nil, err
	}
	return flatbuffers.NewLeagueEncoder().DecodeFromBinary(buf)
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

func (k *Keeper) SaveUnfinishedMatch(ctx context.Context, match datasource.Match) error {
	key := k.MatchKeyUnfinished(match.ID)

	var val = make([]byte, 8)

	fmt.Printf("match.ID: %v\n", match.ID)
	fmt.Printf("val: %v\n", val)
	binary.BigEndian.PutUint64(val, uint64(match.ID))
	return k.storeService.OpenKVStore(ctx).Set(key, val)
}

func (k *Keeper) DeleteUnfinishedMatch(ctx context.Context, matchID int) error {
	key := k.MatchKeyUnfinished(matchID)
	return k.storeService.OpenKVStore(ctx).Delete(key)
}

func (k *Keeper) ListUnfinishedMatches(ctx context.Context) ([]int, error) {
	iterator, err := k.storeService.OpenKVStore(ctx).Iterator(MatchKeyUnfinishedPrefix, nil)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var matchIDs []int
	for ; iterator.Valid(); iterator.Next() {
		val := iterator.Value()
		if len(val) != 8 {
			continue
		}
		id := int(binary.BigEndian.Uint64(val))
		matchIDs = append(matchIDs, id)
	}
	return matchIDs, nil
}
