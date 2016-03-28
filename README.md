![Logo](logo.png)

![License](http://img.shields.io/:license-mit-blue.svg)
[![GoDoc](https://godoc.org/github.com/MaxHalford/gago?status.svg)](https://godoc.org/github.com/MaxHalford/gago)
![Build status](https://api.travis-ci.org/MaxHalford/gago.svg?branch=master)
[![Coverage status](https://coveralls.io/repos/github/MaxHalford/gago/badge.svg?branch=master)](https://coveralls.io/github/MaxHalford/gago?branch=master)
[![Report card](https://img.shields.io/badge/go_report-A+-brightgreen.svg)](https://goreportcard.com/report/github.com/MaxHalford/gago)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/sindresorhus/awesome)

`gago` is a framework for running genetic algorithms. It is written in Go. 

In its most basic form, a [genetic algorithm](https://www.wikiwand.com/en/Genetic_algorithm) solves a mathematically posed problem by doing the following:

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
- Multi-population GAs add another layer in-between called ***demes***.
- Demes exchange individuals through a process known as ***migration***.

## Usage

The following code shows a basic usage of `gago`.

```go
package main

import (
	"math"

	"github.com/MaxHalford/gago"
)

// Sphere function minimum is 0
func sphere(X []float64) float64 {
	sum := 0.0
	for _, x := range X {
		sum += math.Pow(x, 2)
	}
	return sum
}

func main() {
	// Instantiate a population
	ga := gago.Float
	// Wrap the function
	ga.Ff = gago.FloatFunction{sphere}
	// Initialize the genetic algorithm with two variables per individual
	ga.Initialize(2)
	// Enhance the individuals
	for i := 0; i < 20; i++ {
		ga.Enhance()
	}
	// Display the best individual
	ga.Best.Display()
}
```

The user has a function `Sphere` he wishes to minimize. He initiates a predefined set of parameters `gago.Float`. He then wraps his function with `gago.FloatFunction{Sphere}` so that `gago` knows what it has to minimize. Finally the user orders `gago` to find an appropriate solution with `ga.Enhance()`. By convention `gago` will try to **minimize** the fitness function. If instead you want to maximize a function `f(x)`, you can minimize `-f(x)` or `1/f(x)`.

## Parameters

To modify the behavior off the GA, you can change the `gago.Population` struct before running `ga.Initialize`. You can either instantiate a new `gago.Population` or use a predefined one from the `configuration.go` file.

| Variable in the code   | Type                      | Description                                                      |
|------------------------|---------------------------|------------------------------------------------------------------|
| `NbDemes`              | `int`                     | Number of demes in the population.                               |
| `NbIndividuals`        | `int`                     | Number of individuals in each deme.                              |
| `NbGenes`              | `int`                     | Number of genes in each individual.                              |
| `Initializer` (struct) | `Initializer` (interface) | Method for initializing a new individual.                        |
| `Selector` (struct)    | `Selector` (interface)    | Method for selecting one individual from a group of individuals. |
| `Crossover` (struct)     | `Crossover` (interface)     | Method for producing a new individual (called the offspring).    |
| `Mutators` (struct)     | `[]Mutator` (slice of interface)     | Method for modifying an individual's genes.                      |
| `Migrator` (struct)    | `Migrator` (interface)    | Method for exchanging individuals between the demes.             |

The `gago.Population` struct also contains a `Best` variable which is of type `Individual`, it represents the best individual overall demes for the current generation. Alternatively the `Demes` variable is a slice containing each deme in the population; the demes are sorted at each generation so that the first individual in the deme is the best individual from that deme.

`gago` is designed to be flexible. You can change every parameter of the algorithm as long as you implement functions that use the correct types as input/output. A good way to start is to look into the source code and see how the methods are implemented, I've made an effort to comment it. If you want to add a new generic operator (initializer, selector, crossover, mutator, migrator), then you can simply copy and paste an existing method into your code and change the logic as you please. All that matters is that you correctly implement the existing interfaces.

If you wish to not use certain genetic operators, you can set them to `nil`. This is available for the `Crossover`, the `Mutator` and the `Migrator` (the other ones are part of minimum requirements).

## Using different types

Some genetic operators target a specific type, these ones are prefixed with the name of the type (`Float`, `String`). The ones that don't have prefixes work with any types, which is down to the way they are implemented. Default configurations are available in `configuration.go`.

You should think of `gago` as a framework for implementing your problems, and not as an all in one solution. It's quite easy to implement your own for exotic problems, for example the [TSP problem](examples/tsp/).

The only requirement for solving a problem is that the problem can be modeled as a function that returns a floating point value, that's it. Because Go is statically typed, you have to provide a [wrapper for the function](fitness.go) and make sure that the genetic operators make sense for your problem.

## Documentation

- [godoc](https://godoc.org/github.com/MaxHalford/gago)
- Each operator (selection, crossover, mutation, migration) is described in its comments.
- [*An introduction to genetic algorithms*](http://www.boente.eti.br/fuzzy/ebook-fuzzy-mitchell.pdf) is quite thorough.
- [*The Multipopulation Genetic Algorithm: Local Selection and Migration*](http://www.pohlheim.com/Papers/mpga_gal95/gal2_1.html) is an easy read.

## Examples

- Check out the [examples/minimization/](examples/math-functions/) folder for basic examples. Test functions were found [here](http://www.sfu.ca/~ssurjano/optimization.html).
- [examples/plot-fitness/](examples/plot-fitness/) is an example of plotting the fitness per generation with [gonum/plot](https://github.com/gonum/plot).
- [examples/curve-fitting/](examples/curve-fitting/) is an attempt to fit a set of points with non-linear polynomial function.
- [examples/tsp/](examples/tsp/) contain examples of solving the Traveling Salesman Problem.

## Roadmap

- Error handling.
- More tests.
- Statistics.
- Benchmarking.
- Compare with other algorithms/libraries.
- Implement/generalize genetic operators.
- More examples.

## Why use gago?

- It's generic, your only constraint is to model your problem and `gago` will do all the hard work for you.
- You can easily add your own genetic operators.
- `gago` implements parallel populations (called "demes") who exchange individuals for better performance.
- `gago` will be well maintained.

## Suggestions

- Please post suggestions/issues in GitHub's issues section.
- You can use the [reddit thread](https://www.reddit.com/r/golang/comments/43oi5j/gago_a_parallel_genetic_algorithm_with_go/) or my [email address](mailto:maxhalford25@gmail.com) for comments/enquiries.
