# Genetic algorithm in Go

![Logo](logo.png)

In it's most basic form, a [genetic algorithm](https://www.wikiwand.com/en/Genetic_algorithm) runs as follows:

1. Generate random solutions.
2. Evaluate the solutions.
3. Sort the solutions by increasing (or decreasing) order.
4. Apply genetic operators (such as [mutation](https://www.wikiwand.com/en/Mutation_(genetic_algorithm)) and [crossover](http://www.wikiwand.com/en/Crossover_(genetic_algorithm))).
5. Repeat from step 2 until satisfied.

Genetic algorithms can be applied to many problems, the only variable being the problem itself. Indeed, the underlying structure does not have to change between problems. With this in mind, `gago` has been built to be reusable. What's more, `gago` is a [multi-population genetic algorithm](http://www.pohlheim.com/Papers/mpga_gal95/gal2_1.html), in that sense it performs better than a traditional genetic algorithm.

Genetic algorithms are notorious for being [embarrassingly parallel](http://www.wikiwand.com/en/Embarrassingly_parallel). Indeed, most calculations can be run in parallel because they only affect one individual. Luckily Go provides good support for parallelism. As some gophers may know, the `math/rand` module can be problematic because there is a global lock the random number generator, the problem is described in this [stackoverflow post](http://stackoverflow.com/questions/14298523/why-does-adding-concurrency-slow-down-this-golang-code). This can be circumvented by providing each deme with it's own generator.

## Terminology

The terms surrounding genetic algorithms (GAs) are roughly analogous to those found in biology.

- GAs are intended to optimize a function called the ***fitness function***.
- Candidate solutions to the optimization problem are called ***individuals***.
- Each individual has ***genes*** to which a ***fitness*** is assigned.
- The number of genes is simply the number of variables defined by the problem.
- Individuals are sorted based on their fitness towards a problem.
***crossover***.
- Genes can be randomly modified through ***mutation***.
- Classically, individuals are contained in a structure called a ***population***.
- Parallel GAs add another layer in-between called ***demes***.

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
- `ga.Initialize` will populate the demes with individuals.
- The GA will try to **minimize** the fitness function. If instead you want to maximize a function `f(x)`, you can minimize `-f(x)` or `1/f(x)`.

## Parameters

To modify the behavior off the GA, you can change the `gago.Population` struct before running `ga.Initialize`. You can either instantiate a new `gago.Population` or use a predefined one from the `configuration.go` file.

| Variable in the code | Type                                       | Description                                                      |
|----------------------|--------------------------------------------|------------------------------------------------------------------|
| `NbDemes`            | `int`                                      | Number of demes in the population.                               |
| `NbIndividuals`      | `int`                                      | Number of individuals in each deme.                              |
| `NbGenes`            | `int`                                      | Number of genes in each individual.                              |
| `Ff`                 | `func([]float64) float64`                  | Fitness function the GA has to minimize.                         |
| `Boundary`           |                                            | Boundary when generating initial genes.                          |
| `Selection`          | `func(Individuals, *rand.Rand) Individual` | Method for selecting one individual from a group of individuals. |
| `Crossover`          | `func(Individuals, *rand.Rand) Individual` | Method for producing a new individual.                           |
| `CSize`              | `int`                                      | Number of individuals that are chosen for crossover.             |
| `Mutate`             | `func(*Individual, float64, *rand.Rand)`   | Method for modifying an individual's genes.                      |
| `MRate`              | `MRate`                                    | Rate at which genes mutate.                                      |

`gago` is very flexible. You can change every parameter of the algorithm as long as you implement functions that use the correct types as input/output. A good way to start is to look into the source code and see how the methods are implemented, I've made an effort to comment it.

## Examples

- Check out the [examples/minimization/](examples/minimization/) folder for basic examples. Test functions were found [here](http://www.sfu.ca/~ssurjano/optimization.html).
- [examples/plot-fitness/](examples/plot-fitness/) is an example of plotting the fitness per generation with [gonum/plot](https://github.com/gonum/plot).
- [examples/curve-fitting](examples/curve-fitting/) is an attempt to fit a set of points with non-linear polynomial function.

## Roadmap

- Discrete functions.
- Error handling.
- Testing.
- Benchmarking.
- Comparison with other algorithms/libraries.
