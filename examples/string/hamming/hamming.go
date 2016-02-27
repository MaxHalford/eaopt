package main

import (
	"strings"

	"github.com/MaxHalford/gago"
)

// Target string the GA has to guess
var target = strings.Split("hello", "")

func hamming(guess []string) float64 {
	var score = 0
	// For each character in the guess string
	for i := range guess {
		// Check if it matches the target string
		if guess[i] != target[i] {
			score++
		}
	}
	return float64(score)
}

func main() {
	// Instantiate a population
	ga := gago.String
	// Wrap the evaluation function
	ga.Ff = gago.StringFunction{hamming}
	// Initialize the genetic algorithm with the length of the target
	ga.Initialize(len(target))
	// Enhancement
	for i := 0; i < 100; i++ {
		ga.Best.Display()
		ga.Enhance()
	}
}
