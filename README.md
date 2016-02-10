# Genetic algorithm in Go (gago)

![License](http://img.shields.io/:license-mit-blue.svg)
[![GoDoc](https://godoc.org/github.com/MaxHalford/gago?status.svg)](https://godoc.org/github.com/MaxHalford/gago)

![Logo](logo.png)

In it's most basic form, a [genetic algorithm](https://www.wikiwand.com/en/Genetic_algorithm) solves a mathematically posed problem by doing the following:

1. Generate random solutions.
2. Evaluate the solutions.
3. Sort the solutions by increasing (or decreasing) order.
4. Apply genetic operators (such as [mutation](https://www.wikiwand.com/en/Mutation_(genetic_algorithm)) and [crossover](http://www.wikiwand.com/en/Crossover_(genetic_algorithm))).
5. Repeat from step 2 until satisfied.

Genetic algorithms can be applied to many problems, the only variable being the problem itself. Indeed, the underlying structure does not have to change between problems. With this in mind, `gago` has been built to be reusable. What's more, `gago` is a [multi-population genetic algorithm](http://www.pohlheim.com/Papers/mpga_gal95/gal2_1.html) implementing the *migration model*, in that sense it performs better than a traditional genetic algorithm.

Genetic algorithms are notorious for being [embarrassingly parallel](http://www.wikiwand.com/en/Embarrassingly_parallel). Indeed, most calculations can be run in parallel because they only affect one individual. Luckily Go provides good support for parallelism. As some gophers may know, the `math/rand` module can be problematic because there is a global lock the random number generator, the problem is described in this [stackoverflow post](http://stackoverflow.com/questions/14298523/why-does-adding-concurrency-slow-down-this-golang-code). This can be circumvented by providing each deme with it's own generator.

## Terminology

The terms surrounding genetic algorithms (GAs) are roughly analogous to those found in biology.

- GAs are intended to optimize a function called the ***fitness function***.
- Candidate solutions to the optimization problem are called ***individuals***.
- Each individual has ***genes*** to which a ***fitness*** is assigned.
- The number of genes is simply the number of variables defined by the problem.
- Individuals are sorted based on their fitness towards a problem.
- ***Offsprings*** are born by applying ***crossover*** on selected individuals.
- The ***selection*** method is crucial and determines most of the behavior of the algorithm.
- Genes can be randomly modified through ***mutation***.
- Classically, individuals are contained in a structure called a ***population***.
- Multi-population GAs add another layer in-between called ***demes***,

## Usage

The following code shows a basic usage of `gago`.

```go
package main

import (
	"fmt"
	m "math"

	"github.com/MaxHalford/gago"
)

// Sphere function minimum is 0
func Sphere(X []float64) float64 {
	sum := 0.0
	for _, x := range X {
		sum += m.Pow(x, 2)
	}
	return sum
}

func main() {
	// Instantiate a population
	ga := gago.Default
	// Fitness function
	function := Sphere
	// Number of variables the function takes as input
	variables := 2
	// Initialize the genetic algorithm
	ga.Initialize(function, variables)
	// Enhancement
	for i := 0; i < 20; i++ {
		fmt.Println(ga.Best)
		ga.Enhance()
	}
}
```

- `gago.Default` is a preset configuration, this way the parameters described in the following don't have to all be set manually.
- `function` and `variables` have to be predefined. These are only two parameters that change from one problem to another.
- `ga.Initialize` will populate the demes with individuals based on the chosen parameters.
- The GA will try to **minimize** the fitness function. If instead you want to maximize a function `f(x)`, you can minimize `-f(x)` or `1/f(x)`.

## Parameters

To modify the behavior off the GA, you can change the `gago.Population` struct before running `ga.Initialize`. You can either instantiate a new `gago.Population` or use a predefined one from the `configuration.go` file.

| Variable in the code | Type                                                                                  | Description                                                      |
|----------------------|---------------------------------------------------------------------------------------|------------------------------------------------------------------|
| `NbDemes`            | `int`                                                                                 | Number of demes in the population.                               |
| `NbIndividuals`      | `int`                                                                                 | Number of individuals in each deme.                              |
| `NbGenes`            | `int`                                                                                 | Number of genes in each individual.                              |
| `Ff`                 | `func([]float64) float64`                                                             | Fitness function the GA has to minimize.                         |
| `Boundary`           | `float64`                                                                             | Boundary when generating initial genes.                          |
| `SelMethod`          | `func(Individuals, *rand.Rand) Individual`                                            | Method for selecting one individual from a group of individuals. |
| `CrossMethod`        | `func(Individuals, *rand.Rand) Individual`                                            | Method for producing a new individual (called the offspring).    |
| `CrossSize`          | `int`                                                                                 | Number of individuals that are chosen for crossover.             |
| `MutMethod`          | `func(indi *Individual, rate float64, intensity float64, generator *rand.Rand)`       | Method for modifying an individual's genes.                      |
| `MutRate`            | `float64`                                                                             | Rate at which genes mutate.                                      |
| `MutMethod`          | `float64`                                                                             | Intensity at which genes mutate.                                 |
| `MigMethod`		   | `func(demes []Deme) []Deme`														   | Protocol for exchanging individuals between the demes.			  |

`gago` is very flexible. You can change every parameter of the algorithm as long as you implement functions that use the correct types as input/output. A good way to start is to look into the source code and see how the methods are implemented, I've made an effort to comment it.

## Display

- The `String()` methods of each structure have been redefined. For example you can do `fmt.Println(individual)` and the result will be prettily displayed. These display methods are available in the `display.go` script.
- Later on the goal will be to allow exporting information via CSV files etc.

## Documentation

- [godoc](https://godoc.org/github.com/MaxHalford/gago)
- Each operator (selection, crossover, mutation, migration) is described in it's comments.
- [An introduction to genetic algorithms](http://www.boente.eti.br/fuzzy/ebook-fuzzy-mitchell.pdf) is quite thorough.
- [The Multipopulation Genetic Algorithm: Local Selection and Migration](http://www.pohlheim.com/Papers/mpga_gal95/gal2_1.html) is an easy read.

## Examples

- Check out the [examples/minimization/](examples/minimization/) folder for basic examples. Test functions were found [here](http://www.sfu.ca/~ssurjano/optimization.html).
- [examples/plot-fitness/](examples/plot-fitness/) is an example of plotting the fitness per generation with [gonum/plot](https://github.com/gonum/plot).
- [examples/curve-fitting/](examples/curve-fitting/) is an attempt to fit a set of points with non-linear polynomial function.

## Roadmap

- Discrete functions.
- Error handling.
- Testing.
- Benchmarking.
- Comparison with other algorithms/libraries.
- Implement more genetic operators.

## Comments

- Please post suggestions/issues in GitHub's issues section.
- You can use the [reddit thread](https://www.reddit.com/r/golang/comments/43oi5j/gago_a_parallel_genetic_algorithm_with_go/) or my [email address](mailto:maxhalford25@gmail.com) for comments/enquiries.
- I'm quite happy with the syntax and the naming in general, however things are not set in stone and some stuff may change to incorporate more functionalities.
- Genetic algorithms are a deep academic interest of mine, I am very motivated to maintain `gago` and implement state-of-the-art methods.
- As far as I know the `GOMAXPROCS` from the `runtime` library defaults to the number of available thread, hence we haven't set it in the source code.

## Change log

| Date       | Description                                                                                                                                                                                         |
|------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 01/02/2016 | First commit.                                                                                                                                                                                       |
| 02/02/2016 | Based on the apparent popularity of  `gago`, I made some decisions to make it more flexible and readable. Essentially some names have changed and display functions have started to be implemented. |
| 03/02/2016 | The first migration method has been implemented, the documentation has been updated accordingly. Most methods have been capitalized for `godoc` purposes.                                           |
