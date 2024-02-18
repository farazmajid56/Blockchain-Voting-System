package main

import (
	"crypto/sha256"
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
	list := Blockchain[len(Blockchain)-1].Votes
	for i := 0; i < len(list); i++ {
		if voterID == list[i].VoterID {
			return true
		}
	}
	return false
}

// RegisterVoter adds a new voter to the system.
func RegisterVoter(voterID int) {
	fmt.Printf("Voter %d registered.\n", voterID)
	// TODO: Check for duplicates
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

	var votes []Vote
	votes = append(votes, vote)

	lastBlock := Blockchain[len(Blockchain)-1]

	newBlock := Block{
		PrevHash: lastBlock.CurrentHash,
		Votes:    votes, //append(lastBlock.Votes, vote),
	}
	newBlock.CurrentHash = calculateHash(newBlock, vote)

	Blockchain = append(Blockchain, newBlock)

	fmt.Printf("Vote cast by Voter %d for %s is recorded.\n", voterID, candidate)
}

// calculateHash calculates the hash of a block.
func calculateHash(block Block, vote Vote) string {
	data := fmt.Sprintf("%v%v%v", block.PrevHash, block.Votes, vote)
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// CalculateElectionResults calculates and displays the winner of the election.
func CalculateElectionResults() {
	fmt.Println("\nElection Results:")
	var winner string
	maxVotes := 0

	// Write your logic for calculating the winner of the election
	var Results = make(map[string]int)

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
