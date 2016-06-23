<div align="center">
  <!-- Logo -->
   <img src="docs/img/logo.png" alt="logo"/>
</div>

<div align="center">
  <!-- License -->
  <a href="https://opensource.org/licenses/MIT">
    <img src="http://img.shields.io/:license-mit-ff69b4.svg?style=flat-square" alt="logo"/>
  </a>
  <!-- godoc -->
  <a href="https://godoc.org/github.com/MaxHalford/gago">
    <img src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square" alt="godoc" />
  </a>
  <!-- readthedocs -->
  <a href="http://gago.readthedocs.io/">
    <img src="https://img.shields.io/badge/docs-latest-blue.svg?style=flat-square" alt="readthedocs" />
  </a>
  <!-- Build status -->
  <a href="https://img.shields.io/travis/MaxHalford/gago">
    <img src="https://img.shields.io/travis/MaxHalford/gago.svg?style=flat-square" alt="build_status" />
  </a>
  <!-- Test coverage -->
  <a href="https://coveralls.io/github/MaxHalford/gago?branch=master">
    <img src="https://coveralls.io/repos/github/MaxHalford/gago/badge.svg?branch=master&style=flat-square" alt="test_coverage" />
  </a>
  <!-- Go report card -->
  <a href="https://goreportcard.com/report/github.com/MaxHalford/gago">
    <img src="https://goreportcard.com/badge/github.com/MaxHalford/gago?style=flat-square" alt="go_report_card" />
  </a>
  <!-- Dependencies -->
  <a href="https://godoc.org/github.com/MaxHalford/gago?imports">
    <img src="https://img.shields.io/badge/dependencies-none-brightgreen.svg?style=flat-square" alt="dependencies" />
  </a>
  <!-- Go awesome -->
  <a href="https://github.com/sindresorhus/awesome">
    <img src="https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg" alt="go_awesome" />
  </a>
</div>

<br/>

<div align="center"><code>gago</code> is a <b>framework</b> written in <b>Go</b> for running <b>genetic algorithms</b></div>

<br/>

# Introduction

In a nutshell, a [genetic algorithm](https://www.wikiwand.com/en/Genetic_algorithm) (GA) solves an optimization problem by doing the following:

1. Generate random solutions.
2. Evaluate the solutions.
3. Sort the solutions according to their evaluation score.
4. Apply genetic operators following a model.
5. Repeat from step 2 until satisfied.

Genetic algorithms can be applied to many problems, the only variable being the problem itself. Indeed, the underlying structure does not have to change between problems. With this in mind, `gago` has been built to be applicable to many problems without having to re-write boilerplate code.

The following flowchart shows the steps the algorithms takes. As can be seen only the search for the best individual and the migration can't be parallelized.

<br/>
<div align="center">
	<img src="docs/img/flowchart.png" alt="Flowchart"/>
</div>
<br/>

This description is voluntarily vague as to how the genetic operators are applied. It's important to understand that there isn't a single way of applying genetic algorithms. For example some people believe that crossover is useless and use mutation for generating new individuals. Genetic operators are applied following a *model*, a fact that is often omitted in introductions to genetic algorithms.


## Terminology




## Philosophy

> TLDR: gago gives you the freedom to tweak every single part of a genetic algorithm.


## Usage

The following code shows a basic usage of `gago`.

```go
package main

import (
	"fmt"
	"math"

	"github.com/MaxHalford/gago/presets"
)

// Sphere function minimum is 0 reached in (0, ..., 0).
// Any search domain is fine.
func sphere(X []float64) float64 {
	sum := 0.0
	for _, x := range X {
		sum += math.Pow(x, 2)
	}
	return sum
}

func main() {
	// Instantiate a GA with 2 variables and the fitness function
	var ga = presets.Float(2, sphere)
	ga.Initialize()
	// Enhancement
	for i := 0; i < 10; i++ {
		ga.Enhance()
	}
	// Display the best obtained solution
	fmt.Printf("The best obtained solution is %f\n", ga.Best.Fitness)
}
```

The user has a function `Sphere` he wishes to minimize. He initiates a predefined set of parameters `presets.Float(2, sphere)`. Finally the user orders `gago` to find an appropriate solution with `ga.Enhance()`. By convention `gago` will try to **minimize** the fitness function. If instead you want to maximize a function `f(x)`, you can minimize `-f(x)`.

The nice thing thing is that the GA only requires a function and a number of variables, hence it can be used for minimizing any function as long the function outputs a floating point number (`float64` in Go). Because Go is a statically typed language and doesn't provide explicit generic types, the input type of the function has be handled. This is done using interfaces and is [documented further down](#using-different-types).


## Features

- Possibility to run many populations in parallel.
- Custom genetic operators are easy to implement, as described in the [contribution document](docs/CONTRIBUTING.md).
- Possibility to apply genetic operators under custom models.
- A modular approach makes it easy to switch GA parameters.
- Speciation operators to cluster individuals into similar groups, providing more efficient crossovers.


## Suggestions

- Please post suggestions/issues in GitHub's issues section.
- Check out the [contribution documentation](CONTRIBUTING.md) if you want to want to participate to `gago`.
- You can use the [reddit thread](https://www.reddit.com/r/golang/comments/43oi5j/gago_a_parallel_genetic_algorithm_with_go/) or my [email address](mailto:maxhalford25@gmail.com) for comments/enquiries.


# Documentation

Don't forget to also check out [godoc](https://godoc.org/github.com/MaxHalford/gago) for a full list of accessible variables.

## Parameters

To modify the behavior off the GA, you can change the `gago.GA` struct before running `ga.Initialize`. You can either instantiate a new `gago.GA` or use a predefined one from the `configuration.go` file.

| Variable in the code   | Type                      | Description                                                      |
|------------------------|---------------------------|------------------------------------------------------------------|
| `NbPopulations`              | `int`                     | Number of Populations in the GA                               |
| `NbIndividuals`        | `int`                     | Number of individuals in each population                              |
| `Initializer` (struct) | `Initializer` (interface) | Method for initializing a new individual                        |
| `Model` (struct)    | `Model` (interface)    | Model for applying genetic operators |
| `Migrator` (struct)    | `Migrator` (interface)    | Method for exchanging individuals between the populations             |
| `MigFrequency`     | `int`     | Migration frequency                      |

The `gago.GA` struct also contains a `Best` variable which is of type `Individual`, it represents the best individual overall. The `Populations` variable is a slice containing each GA in the GA. The populations are sorted at each generation so that the first individual in each GA is the best individual for that specific GA.

`gago` is designed to be flexible. You can change every parameter of the algorithm as long as you implement functions that use the correct types as input/output. A good way to start is to look into the source code and see how the methods are implemented, I've made an effort to comment each and every one of them. If you want to add a new generic operator (initializer, selector, crossover, mutator, migrator), then you can simply copy and paste an existing method into your code and change the logic as you see fit. All that matters is that you correctly implement the existing interfaces.

If you wish to not use certain genetic operators, you can set them to `nil`. This is available for the `Mutator` and the `Migrator` (the other ones are part of the minimum requirements). Each operator contains an explanatory description that can be consulted in the [documentation](https://godoc.org/github.com/MaxHalford/gago).


## Using different types

Some genetic operators target a specific type, these ones are suffixed with the name of the type (`F` for `float64`, `S` for `string`). The ones that don't have suffixes work with any types, which is down to the way they are implemented.

You should think of `gago` as a framework for implementing your problems, and not as an all in one solution. It's quite easy to implement custom operators for exotic problems, for example the [TSP problem](examples/tsp/).

The only requirement for solving a problem is that the problem itself can be modeled as a function that returns a floating point value. Because Go is statically typed, you have to provide a [wrapper for the function](fitness.go) and make sure that the genetic operators make sense for your problem. The reasoning behing `gago` makes more sense once you start looking at the examples.


## Presets

For conveniency `gago` includes a set of [presets](presets/). For the while these serve the purpose of exemplifying possible configurations. Later on these could include state of the art sets of parameters based on battle testing on common problems. For examples it could be useful to tuning a preset for solving TSP problems.


## Advice

- Wrap multiple mutators into a single `Mutator` `struct` if you wish to apply multiple mutators, for an example see the [TSP preset](presets/tsp.go).
- Don't hesitate to add more populations if you have a multi-core machine, the overhead is very small.
- Consider the fact that most of the computation is for evaluating the fitness function.
- Increasing the number of selected parents (`NbParents`) during selection usually increases the convergence rate (which is not necessarily good, but is sometimes desired).
- Increasing the number of individuals per population (`NbIndividuals`) adds variety to the genetic algorithm, however it is more costly.
- You can access the GA's `duration` attribute or implement your own stopwatch to enhance the GA for a fixed duration.


## Examples

- Check out the [examples/minimization/](examples/math-functions/) folder for basic examples. Test functions were found [here](http://www.sfu.ca/~ssurjano/optimization.html).
- [examples/plot-fitness/](examples/plot-fitness/) is an example of plotting the fitness per generation with [gonum/plot](https://github.com/gonum/plot).
- [examples/curve-fitting/](examples/curve-fitting/) is an attempt to fit a set of points with non-linear polynomial function.
- [examples/tsp/](examples/tsp/) contain examples of solving the Traveling Salesman Problem.

