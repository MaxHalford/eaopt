![Logo](logo.png)

# gago is a genetic algorithm written in Golang

![License](http://img.shields.io/:license-mit-blue.svg)
[![GoDoc](https://godoc.org/github.com/MaxHalford/gago?status.svg)](https://godoc.org/github.com/MaxHalford/gago)

In it's most basic form, a [genetic algorithm](https://www.wikiwand.com/en/Genetic_algorithm) solves a mathematically posed problem by doing the following:

1. Generate random solutions.
2. Evaluate the solutions.
3. Sort the solutions by increasing (or decreasing) order.
4. Apply genetic operators (such as [mutation](https://www.wikiwand.com/en/Mutation_(genetic_algorithm)) and [breeding](http://www.wikiwand.com/en/Breeder_(genetic_algorithm))).
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
- ***Offsprings*** are born by applying ***breeding*** on selected individuals.
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
| `Breeder` (struct)     | `Breeder` (interface)     | Method for producing a new individual (called the offspring).    |
| `Mutator` (struct)     | `Mutator` (interface)     | Method for modifying an individual's genes.                      |
| `Migrator` (struct)    | `Migrator` (interface)    | Method for exchanging individuals between the demes.             |

The `gago.Population` struct also contains a `Best` variable which is of type `Individual`, it represents the best individual overall demes for the current generation. Alternatively the `Demes` variable is a slice containing each deme in the population; the demes are sorted at each generation so that the first individual in the deme is the best individual from that deme.

`gago` is designed to be flexible. You can change every parameter of the algorithm as long as you implement functions that use the correct types as input/output. A good way to start is to look into the source code and see how the methods are implemented, I've made an effort to comment it. If you want to add a new generic operator (initializer, selector, breeding, mutator, migrator), then you can simply copy and paste an existing method into your code and change the logic as you please. All that matters is that you correctly implement the existing interfaces.

If you wish to not use certain genetic operators, you can set them to `nil`. This is available for the `Breeder`, the `Mutator` and the `Migrator`.


## Using different types

It works, documentation is coming.


## Documentation

- [godoc](https://godoc.org/github.com/MaxHalford/gago)
- Each operator (selection, breeding, mutation, migration) is described in it's comments.
- [**An introduction to genetic algorithms**](http://www.boente.eti.br/fuzzy/ebook-fuzzy-mitchell.pdf) is quite thorough.
- [**The Multipopulation Genetic Algorithm: Local Selection and Migration**](http://www.pohlheim.com/Papers/mpga_gal95/gal2_1.html) is an easy read.

## Examples

- Check out the [examples/minimization/](examples/minimization/) folder for basic examples. Test functions were found [here](http://www.sfu.ca/~ssurjano/optimization.html).
- [examples/plot-fitness/](examples/plot-fitness/) is an example of plotting the fitness per generation with [gonum/plot](https://github.com/gonum/plot).
- [examples/curve-fitting/](examples/curve-fitting/) is an attempt to fit a set of points with non-linear polynomial function.

## Roadmap

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
| 13/02/2016 | The genetic operator API got an overhaul. Each operator now implements an interface, this makes things more readable and more flexible.                                                             |
| 14/02/2016 | Any type of fitness function is now accepted as long as the correct interfaces are implemented. Most of the methods have become private to minimize the API.                                        |
