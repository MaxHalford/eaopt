package main

import (
	"fmt"
	"strings"

	"github.com/MaxHalford/gago/presets"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numerals = "0123456789"
	symbols  = " ~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"
	ascii    = alphabet + numerals + symbols
)

var (
	corpus = strings.Split(ascii, "")
	target = strings.Split("Hello world!", "")
)

func hamming(str []string) float64 {
	var errors float64
	for i, s := range str {
		if s != target[i] {
			errors++
		}
	}
	return errors
}

func main() {
	// Create the GA
	var ga = presets.Alignment(len(target), corpus, hamming)
	ga.Initialize()
	// Enhance
	for i := 0; i < 3000; i++ {
		ga.Enhance()
	}
	// Extract the genome of the best individual
	var str = ga.Best.Genome.CastString()
	fmt.Println(strings.Join(str, " "))
}
