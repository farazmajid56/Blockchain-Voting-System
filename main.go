package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// Vote represents a single vote cast by a voter.
type Vote struct {
	VoterID   int
	Candidate string
}

// Block represents a block in the blockchain, containing multiple votes.
type Block struct {
	PrevHash    string
	CurrentHash string // Add CurrentHash field
	Votes       []Vote
}

// Blockchain is a slice of Block elements.
var Blockchain []Block

// Candidates is a map to store candidate vote counts.
var Candidates map[string]int

// Registered Voters // SET
var RegisteredVoters []int

func IsRegisteredVoter(voterID int) bool {
	for i := 0; i < len(RegisteredVoters); i++ {
		if RegisteredVoters[i] == voterID {
			return true
		}
	}
	return false
}

func IsDuplicateVote(voterID int) bool {
	for i := 1; i < len(Blockchain); i++ {
		if voterID == Blockchain[i].Votes[0].VoterID {
			return true
		}
	}
	return false
}

// RegisterVoter adds a new voter to the system.
func RegisterVoter(voterID int) {
	fmt.Printf("Voter %d registered.\n", voterID)
	for _, value := range RegisteredVoters {
		if value == voterID {
			fmt.Printf("Voter %d has already registered.\n", voterID)
			return
		}
	}
	RegisteredVoters = append(RegisteredVoters, voterID)
}

func IsValidCandidate(candidate string) bool {
	_, exists := Candidates[candidate]
	return exists
}

// CastVote allows a registered voter to cast a vote.
func CastVote(voterID int, candidate string) {
	// Check if the voter is registered
	if !IsRegisteredVoter(voterID) {
		fmt.Printf("Invalid voter ID: %d\n", voterID)
		return
	}

	// Check if the voter has already cast a vote
	if IsDuplicateVote(voterID) {
		fmt.Printf("Voter %d has already cast a vote.\n", voterID)
		return
	}

	// Check if the candidate exists
	if !IsValidCandidate(candidate) {
		fmt.Printf("Candidate %s does not exist.\n", candidate)
		return
	}

	// Add the vote to the current block
	vote := Vote{VoterID: voterID, Candidate: candidate}

	lastBlock := Blockchain[len(Blockchain)-1]

	newBlock := Block{
		PrevHash: lastBlock.CurrentHash,
		Votes:    []Vote{vote}, //append(lastBlock.Votes, vote),
	}

	newHash := calculateHash(newBlock, vote)
	if newHash == "error" {
		fmt.Print("Error in conversion to Bytes")
		return
	}

	newBlock.CurrentHash = newHash //calculateHash(newBlock, vote)

	Blockchain = append(Blockchain, newBlock)

	// Update candidate vote count
	Candidates[candidate]++

	fmt.Printf("Vote cast by Voter %d for %s is recorded.\n", voterID, candidate)
}

func ConvertDataToBytes(vote Vote, votes []Vote, prevHash string) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Convert Vote to bytes
	if err := binary.Write(buf, binary.LittleEndian, int64(vote.VoterID)); err != nil {
		return nil, fmt.Errorf("error converting vote.VoterID to bytes: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, []byte(vote.Candidate)); err != nil {
		return nil, fmt.Errorf("error converting vote.Candidate to bytes: %v", err)
	}

	// Convert []Vote to bytes
	for _, v := range votes {
		if err := binary.Write(buf, binary.LittleEndian, int64(v.VoterID)); err != nil {
			return nil, fmt.Errorf("error converting vote.VoterID to bytes: %v", err)
		}
		if err := binary.Write(buf, binary.LittleEndian, []byte(v.Candidate)); err != nil {
			return nil, fmt.Errorf("error converting vote.Candidate to bytes: %v", err)
		}
	}

	// Convert PrevHash to bytes
	if err := binary.Write(buf, binary.LittleEndian, []byte(prevHash)); err != nil {
		return nil, fmt.Errorf("error converting prevHash to bytes: %v", err)
	}

	return buf.Bytes(), nil
}

// calculateHash calculates the hash of a block.
func calculateHash(block Block, vote Vote) string {
	concatenatedData, err := ConvertDataToBytes(vote, block.Votes, block.PrevHash)
	if err != nil {
		return "error"
	}
	hash := sha256.New()
	hash.Write([]byte(concatenatedData))
	return hex.EncodeToString(hash.Sum(nil))
}

// CalculateElectionResults calculates and displays the winner of the election.
func CalculateElectionResults() {
	fmt.Println("\nElection Results:")
	var winner string
	maxVotes := 0

	// Write your logic for calculating the winner of the election
	var Results = make(map[string]int)
	// * Normal Version (For Dummies 🤓)
	// for candidate, votes := range Candidates {
	// 	fmt.Printf("%s: %d votes\n", candidate, votes)
	// 	if votes > maxVotes {
	// 		maxVotes = votes
	// 		winner = candidate
	// 	} else if maxVotes == votes {
	// 		winner = "Tie"
	// 	}
	// }

	// * Blockchain Version Expensive !!!
	// Ignore Genisis Block
	for i := 1; i < len(Blockchain); i++ {
		Results[Blockchain[i].Votes[0].Candidate]++
	}

	for key, value := range Results {
		fmt.Printf("%s: %d votes\n", key, value)
	}

	// Find the winner and check for ties
	var tieCandidates []string // Keeps track of maxVote candidates
	for candidate, votes := range Results {
		if votes > maxVotes {
			maxVotes = votes
			winner = candidate
			tieCandidates = nil
			tieCandidates = append(tieCandidates, candidate)
		} else if votes == maxVotes {
			tieCandidates = append(tieCandidates, candidate)
		}
	}

	if len(tieCandidates) > 1 {
		winner = "Tie"
	}

	if winner != "Tie" {
		fmt.Printf("Winner: %s\n", winner)
	} else {
		fmt.Println("Election resulted in a tie.")
	}
}

func main() {
	// Initialize the blockchain with a genesis block.
	genesisBlock := Block{PrevHash: "", CurrentHash: "", Votes: nil}
	Blockchain = append(Blockchain, genesisBlock)

	// Register candidates
	Candidates = make(map[string]int)
	Candidates["Candidate A"] = 0
	Candidates["Candidate B"] = 0

	// Register voters
	for i := 1; i <= 10; i++ {
		RegisterVoter(i)
	}
	fmt.Print("\n")
	// Simulate voting process
	CastVote(1, "Candidate A")
	CastVote(2, "Candidate B")
	CastVote(3, "Candidate A")
	CastVote(3, "Candidate B") // Attempted Double Voting
	CastVote(4, "Candidate B")
	CastVote(5, "Candidate A")
	CastVote(5, "Candidate A")  // Attempted Double Voting
	CastVote(6, "Candidate B")  // Should Print in case of tie
	CastVote(7, "Candidate C")  // Invalid Candidate ID
	CastVote(11, "Candidate B") // Invalid Voter ID

	// Calculate and display election results
	CalculateElectionResults()

	// Display the blockchain
	fmt.Println("\nBlockchain:")
	for i, block := range Blockchain {
		fmt.Printf("Block %d\n", i)
		fmt.Printf("PrevHash: %s\n", block.PrevHash)
		fmt.Printf("CurrentHash: %s\n", block.CurrentHash)
		fmt.Printf("Votes: %v\n\n", block.Votes)
	}
}
