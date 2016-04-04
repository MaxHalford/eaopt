// Package gago has a convention for naming genetic operators. The name begins
// with a letter or a short sequence of letters to specify which kind of
// operator it is:
//
// - `C`: crossover
// - `Mut`: mutator
// - `Mig`: migrator
// - `S`: selector
//
// Then comes the second part of the name which indicates on what kind of
// genomes the operator works:
//
// - `F`: `float64`
// - `S`: `string`
// - No letter means the operator works on any kind of genome, regardless of the
// underlying type.
//
// Finally the name of the operator ends with a word to indicate what it does.
package gago
