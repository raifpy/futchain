package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	TeamKey   = []byte("team")
	MatchKey  = []byte("match")
	LeagueKey = []byte("league")

	MatchKeyUnfinishedPrefix = []byte("match_unfinished")
)

func (k *Keeper) TeamKey(id int) []byte {

	return append(TeamKey, sdk.Uint64ToBigEndian(uint64(id))...)
}

func (k *Keeper) MatchKey(id int) []byte {
	return append(MatchKey, sdk.Uint64ToBigEndian(uint64(id))...)
}

func (k *Keeper) LeagueKey(id int) []byte {
	return append(LeagueKey, sdk.Uint64ToBigEndian(uint64(id))...)
}

// list only unfinished matchs prefix; Once the match finishs, we should delete it from the store.
func (k *Keeper) MatchKeyUnfinished(id int) []byte {
	return append(MatchKeyUnfinishedPrefix, sdk.Uint64ToBigEndian(uint64(id))...)
}
