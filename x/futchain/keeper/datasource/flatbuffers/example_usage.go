package flatbuffers

import (
	"fmt"
	"time"

	"github.com/raifpy/futchain/x/futchain/keeper/datasource"
)

// ExampleUsage demonstrates how to use the FlatBuffers encoder for blockchain storage
func ExampleUsage() {
	// Create a sample team
	team := &datasource.Team{
		ID:       123,
		Score:    2,
		Name:     "Arsenal",
		LongName: "Arsenal Football Club",
	}

	// Initialize encoders
	teamEncoder := NewTeamEncoder()
	universalEncoder := NewUniversalEncoder()

	// Example 1: Basic encoding/decoding
	fmt.Println("=== Basic Team Encoding ===")

	// Encode to binary
	binaryData, err := teamEncoder.EncodeToBinary(team)
	if err != nil {
		fmt.Printf("Error encoding team: %v\n", err)
		return
	}
	fmt.Printf("Encoded size: %d bytes\n", len(binaryData))

	// Decode from binary
	decodedTeam, err := teamEncoder.DecodeFromBinary(binaryData)
	if err != nil {
		fmt.Printf("Error decoding team: %v\n", err)
		return
	}
	fmt.Printf("Decoded team: %+v\n", decodedTeam)

	// Example 2: Hex encoding for debugging/logging
	fmt.Println("\n=== Hex Encoding for Debugging ===")

	hexStr, err := teamEncoder.EncodeToHex(team)
	if err != nil {
		fmt.Printf("Error encoding to hex: %v\n", err)
		return
	}
	fmt.Printf("Hex representation: %s\n", hexStr)

	// Example 3: Size calculation before encoding
	fmt.Println("\n=== Size Calculation ===")

	size, err := teamEncoder.GetSize(team)
	if err != nil {
		fmt.Printf("Error calculating size: %v\n", err)
		return
	}
	fmt.Printf("Predicted size: %d bytes\n", size)

	// Example 4: Complex data structure (Match)
	fmt.Println("\n=== Match Encoding ===")

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

	status := &datasource.Status{
		UtcTime:      time.Now(),
		PeriodLength: 90,
		Started:      true,
		Cancelled:    false,
		Finished:     false,
	}

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

	matchEncoder := NewMatchEncoder()
	matchData, err := matchEncoder.EncodeToBinary(match)
	if err != nil {
		fmt.Printf("Error encoding match: %v\n", err)
		return
	}
	fmt.Printf("Match encoded size: %d bytes\n", len(matchData))

	// Example 5: Using universal encoder
	fmt.Println("\n=== Universal Encoder ===")

	// Encode using universal encoder
	universalData, err := universalEncoder.EncodeTeam(team)
	if err != nil {
		fmt.Printf("Error with universal encoder: %v\n", err)
		return
	}
	fmt.Printf("Universal encoder size: %d bytes\n", len(universalData))

	// Example 6: Blockchain storage simulation
	fmt.Println("\n=== Blockchain Storage Simulation ===")

	// Simulate storing in blockchain (key-value store)
	storageKey := fmt.Sprintf("team_%d", team.ID)
	storageValue := binaryData

	fmt.Printf("Storage Key: %s\n", storageKey)
	fmt.Printf("Storage Value Size: %d bytes\n", len(storageValue))

	// Simulate retrieving from blockchain
	retrievedTeam, err := teamEncoder.DecodeFromBinary(storageValue)
	if err != nil {
		fmt.Printf("Error retrieving from storage: %v\n", err)
		return
	}
	fmt.Printf("Retrieved team: %+v\n", retrievedTeam)
}

// ExampleBlockchainIntegration shows how to integrate with Cosmos SDK store
func ExampleBlockchainIntegration() {
	fmt.Println("\n=== Blockchain Integration Example ===")

	// This is a conceptual example of how you might integrate with Cosmos SDK
	// In practice, you would use the actual store interfaces from your keeper

	team := &datasource.Team{
		ID:       456,
		Score:    3,
		Name:     "Manchester United",
		LongName: "Manchester United Football Club",
	}

	encoder := NewTeamEncoder()

	// Encode for storage
	data, err := encoder.EncodeToBinary(team)
	if err != nil {
		fmt.Printf("Error encoding: %v\n", err)
		return
	}

	// In a real implementation, you would do something like:
	// store.Set(ctx, []byte(fmt.Sprintf("team_%d", team.ID)), data)
	fmt.Printf("Would store in blockchain: key=team_%d, value_size=%d bytes\n", team.ID, len(data))

	// For retrieval:
	// data := store.Get(ctx, []byte(fmt.Sprintf("team_%d", team.ID)))
	// team, err := encoder.DecodeFromBinary(data)

	fmt.Printf("Team data ready for blockchain storage\n")
}
