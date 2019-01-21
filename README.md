<div align="center">
  <!-- Logo -->
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vRzNhZcghWl-xRaynEAQn4moZeyh6pmD8helS099eJ7V_TRGM9dE6e0T5IF9bOG7S4K62CKZtSUpans/pub?w=477&h=307" alt="logo"/>
</div>

<div align="center">
  <!-- Awesome Go -->
  <a href="https://github.com/avelino/awesome-go">
    <img src="https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg" alt="awesome_go" />
  </a>
  <!-- Awesome Machine Learning -->
  <a href="https://github.com/josephmisiti/awesome-machine-learning">
    <img src="https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg" alt="awesome_ml" />
  </a>
</div>

<br/>

<div align="center">
  <!-- godoc -->
  <a href="https://godoc.org/github.com/MaxHalford/eaopt">
    <img src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square" alt="godoc" />
  </a>
  <!-- Build status -->
  <a href="https://travis-ci.org/MaxHalford/eaopt">
    <img src="https://img.shields.io/travis/MaxHalford/eaopt/master.svg?style=flat-square" alt="build_status" />
  </a>
  <!-- Test coverage -->
  <a href="https://coveralls.io/github/MaxHalford/eaopt?branch=master">
    <img src="https://coveralls.io/repos/github/MaxHalford/eaopt/badge.svg?branch=master&style=flat-square" alt="test_coverage" />
  </a>
  <!-- Go report card -->
  <a href="https://goreportcard.com/report/github.com/MaxHalford/eaopt">
    <img src="https://goreportcard.com/badge/github.com/MaxHalford/eaopt?style=flat-square" alt="go_report_card" />
  </a>
  <!-- Code Climate -->
  <a href="https://codeclimate.com/github/MaxHalford/eaopt">
    <img src="https://codeclimate.com/github/MaxHalford/eaopt/badges/gpa.svg" alt="Code Climate" />
  </a>
  <!-- License -->
  <a href="https://opensource.org/licenses/MIT">
    <img src="http://img.shields.io/:license-mit-ff69b4.svg?style=flat-square" alt="license"/>
  </a>
</div>

<br/>

<div align="center">eaopt is an evolutionary optimization library</div>

<br/>

**Table of Contents**
  - [Changelog](#changelog)
  - [Example](#example)
  - [Background](#background)
  - [Features](#features)
  - [Usage](#usage)
    - [General advice](#general-advice)
    - [Genetic algorithms](#genetic-algorithms)
      - [Overview](#overview)
      - [Implementing the Genome interface](#implementing-the-genome-interface)
      - [Instantiating the GA struct](#instantiating-the-ga-struct)
      - [Calling the Minimize method](#calling-the-minimize-method)
      - [Using the Slice interface](#using-the-slice-interface)
      - [Models](#models)
      - [Multiple populations and migration](#multiple-populations-and-migration)
      - [Speciation](#speciation)
      - [Logging population statistics](#logging-population-statistics)
    - [Particle swarm optimization](#particle-swarm-optimization)
    - [Differential evolution](#differential-evolution)
    - [OpenAI evolution strategy](#openai-evolution-strategy)
  - [A note on parallelism](#a-note-on-parallelism)
  - [FAQ](#faq)
  - [Dependencies](#dependencies)
  - [License](#license)

## Changelog

- **11/11/18**: a simple version of [OpenAI's evolution strategy](https://blog.openai.com/evolution-strategies/) has been implemented, it's called `OES`.
- **02/08/18**: gago has now become eaopt. You can still everything you could do before but the scope is now larger than genetic algorithms. The goal is to implement many more evolutionary optimization algorithms on top of the existing codebase.

## Example

The following example attempts to minimize the [Drop-Wave function](https://www.sfu.ca/~ssurjano/drop.html) using a genetic algorithm. The Drop-Wave function is known to have a minimum value of -1 when each of it's arguments is equal to 0.

<div align="center">
  <img src="https://github.com/MaxHalford/eaopt-examples/blob/master/drop_wave/chart.png" alt="drop_wave_chart" />
  <img src="https://github.com/MaxHalford/eaopt-examples/blob/master/drop_wave/function.png" alt="drop_wave_function" />
</div>

```go
package main

import (
    "fmt"
    m "math"
    "math/rand"

    "github.com/MaxHalford/eaopt"
)

// A Vector contains float64s.
type Vector []float64

// Evaluate a Vector with the Drop-Wave function which takes two variables as
// input and reaches a minimum of -1 in (0, 0). The function is simple so there
// isn't any error handling to do.
func (X Vector) Evaluate() (float64, error) {
    var (
        numerator   = 1 + m.Cos(12*m.Sqrt(m.Pow(X[0], 2)+m.Pow(X[1], 2)))
        denominator = 0.5*(m.Pow(X[0], 2)+m.Pow(X[1], 2)) + 2
    )
    return -numerator / denominator, nil
}

// Mutate a Vector by resampling each element from a normal distribution with
// probability 0.8.
func (X Vector) Mutate(rng *rand.Rand) {
    eaopt.MutNormalFloat64(X, 0.8, rng)
}

// Crossover a Vector with another Vector by applying uniform crossover.
func (X Vector) Crossover(Y eaopt.Genome, rng *rand.Rand) {
    eaopt.CrossUniformFloat64(X, Y.(Vector), rng)
}

// Clone a Vector to produce a new one that points to a different slice.
func (X Vector) Clone() eaopt.Genome {
    var Y = make(Vector, len(X))
    copy(Y, X)
    return Y
}

// VectorFactory returns a random vector by generating 2 values uniformally
// distributed between -10 and 10.
func VectorFactory(rng *rand.Rand) eaopt.Genome {
    return Vector(eaopt.InitUnifFloat64(2, -10, 10, rng))
}

func main() {
    // Instantiate a GA with a GAConfig
    var ga, err = eapot.NewDefaultGAConfig().NewGA()
    if err != nil {
        fmt.Println(err)
        return
    }

    // Set the number of generations to run for
    ga.NGenerations = 10

    // Add a custom print function to track progress
    ga.Callback = func(ga *eapot.GA) {
        fmt.Printf("Best fitness at generation %d: %f\n", ga.Generations, ga.HallOfFame[0].Fitness)
    }

    // Find the minimum
    err = ga.Minimize(VectorFactory)
    if err != nil {
        fmt.Println(err)
        return
    }
}

```

```sh
>>> Best fitness at generation 0: -0.550982
>>> Best fitness at generation 1: -0.924220
>>> Best fitness at generation 2: -0.987282
>>> Best fitness at generation 3: -0.987282
>>> Best fitness at generation 4: -0.987282
>>> Best fitness at generation 5: -0.987282
>>> Best fitness at generation 6: -0.987282
>>> Best fitness at generation 7: -0.997961
>>> Best fitness at generation 8: -0.999954
>>> Best fitness at generation 9: -0.999995
>>> Best fitness at generation 10: -0.999999
```

All the examples can be found [in this repository](https://github.com/MaxHalford/eaopt-examples).

## Background

Evolutionary optimization algorithms are a subdomain of evolutionary computation. Their goal is to minimize/maximize a function without using any gradient information (usually because there isn't any gradient available). They share the common property of exploring the search space by breeding, mutating, evaluating, and sorting so-called *individuals*. Most evolutionary algorithms are designed to handle real valued functions, however in practice they are commonly used for handling more exotic problems. For example genetic algorithms can be used to find the optimal structure of a neural network.

eaopt provides implementations for various evolutionary optimization algorithms. Implementation-wise, the idea is that most (if not all) of said algorithms can be written as special cases of a genetic algorithm. Indeed this is made possible by using a generic definition of a genetic algorithm by allowing the mutation, crossover, selection, and replacement procedures to be modified at will. The `GA` struct is thus the most flexible struct of eaopt, the other algorithms are written on top of it. If you don't find any algorithm that suits your need then you can easily write your own operators (as is done in most of the [examples](https://github.com/MaxHalford/eaopt-examples)).

## Features

- Different evolutionary algorithms are available with a consistent API
- You can practically do anything by using the `GA` struct
- Speciation and migration procedures are available
- Common genetic operators (mutation, crossover, selection, migration, speciation) are already implemented
- Function evaluation can be done in parallel if your function is costly

## Usage

### General advice

- Evolutionary algorithms are usually designed for solving specific kinds of problems. Take a look at the `Minimize` function of each method to get an idea of what type of function it can optimize.
- Use the associated constructor function of each method to initialize it. For example use the `NewPSO` function instead of instantiating the `PSO` struct yourself. Along with making your life easier, these functions provide the added benefit of checking for parameter input errors.
- If you're going to use the `GA` struct then be aware that some evolutionary operators are already implemented in eaopt (you don't necessarily have to reinvent the wheel).
- Don't feel overwhelmed by the fact that algorithms are implemented as special cases of genetic algorithms. It doesn't matter if you just want to get things done, it just makes things easier under the hood.

### Genetic algorithms

#### Overview

Genetic algorithms are the backbone of eaopt. Most of the other algorithms available in eaopt are implemented as special cases of GAs. A GA isn't an algorithm per say, but rather a blueprint which can be used to optimize any kind of problem.

In a nutshell, a GA solves an optimization problem by doing the following:

1. Generate random solutions to a problem.
2. Assign a fitness to each solutions.
3. Check if a new best solution has been found.
4. Apply genetic operators following a given evolutionary model.
5. Repeat from step 2 until the stopping criterion is satisfied.

This description is voluntarily vague. It is up to the user to define the problem and the genetic operators to use. Different categories of genetic operators exist:

- Mutation operators modify an existing solution.
- Crossover operators generate a new solution by combining two or more existing ones.
- Selection operators selects individuals that are to be evolved.
- Migration swaps individuals between populations.
- Speciation clusters individuals into subpopulations.

Popular stopping criteria include

- a fixed number of generations,
- a fixed duration,
- an indicator that the population is stagnating.

Genetic algorithms can be used via the `GA` struct. The necessary steps for using the GA struct are

1. Implement the `Genome` interface to model your problem
2. Instantiate a `GA` struct (preferably via the `GAConfig` struct)
3. Call the GA's `Minimize` function and check the `HallOfFame` field


#### Implementing the Genome interface

To use the `GA` struct you first have to implement the `Genome` interface, which is used to define the logic that is specific to your problem (logic that eaopt doesn't know about). For example this is where you will define an `Evaluate()` method for evaluating a particular problem. The `GA` struct contains context-agnostic information. For example this is where you can choose the number of individuals in a population (which is a separate concern from your particular problem). Apart from a good design pattern, decoupling the problem definition from the optimization through the `Genome` interface means that eaopt can be used to optimize *any* kind of problem.

Let's have a look at the `Genome` interface.

```go
type Genome interface {
    Evaluate() (float64, error)
    Mutate(rng *rand.Rand)
    Crossover(genome Genome, rng *rand.Rand)
    Clone() Genome
}
```

The `Evaluate()` method returns the fitness of a genome. The sweet thing is that you can do whatever you want in this method. Your struct that implements the interface doesn't necessarily have to be a slice. The `Evaluate()` method is *your* problem to deal with. eaopt only needs it's output to be able to function. You can also return an `error` which eaopt will catch and return when calling `ga.Initialize()` and `ga.Evolve()`.

The `Mutate(rng *rand.Rand)` method is where you can modify an existing solution by tinkering with it's variables. The way in which you should mutate a solution essentially boils down to your particular problem. eaopt provides some common mutation methods that you can use instead of reinventing the wheel -- this is what is being done in most of the [examples](https://github.com/MaxHalford/eaopt-examples).

The `Crossover(genome Genome, rng *rand.Rand)` method combines two individuals. The important thing to notice is that the type of first argument differs from the struct calling the method. Indeed the first argument is a `Genome` that has to be casted into your struct before being able to apply a crossover operator. This is due to the fact that Go doesn't provide generics out of the box; it's easier to convince yourself by checking out the examples.

The `Clone()` method is there to produce independent copies of the struct you want to evolve. This is necessary for internal reasons and ensures that pointer fields are not pointing to identical memory addresses. Usually this is not too difficult implement; you just have to make sure that the clones you produce are not shallow copies of the genome that is being cloned. This is also fairly easy to unit test.

Once you have implemented the `Genome` interface you have provided eaopt with all the information it couldn't guess for you.

#### Instantiate the GA struct

You can now instantiate a `GA` and use it to find an optimal solution to your problem. The `GA` struct has a lot of fields, hence the recommended way is to use the `GAConfig` struct and call it's `NewGA` method.

Let's have a look at the `GAConfig` struct.

```go
type GAConfig struct {
    // Required fields
    NPops        uint
    PopSize      uint
    NGenerations uint
    HofSize      uint
    Model        Model

    // Optional fields
    ParallelEval bool // Whether to evaluate Individuals in parallel or not
    Migrator     Migrator
    MigFrequency uint // Frequency at which migrations occur
    Speciator    Speciator
    Logger       *log.Logger
    Callback     func(ga *GA)
	EarlyStop    func(ga *GA) bool
    RNG          *rand.Rand
}
```

- Required fields
  - `NPops` determines the number of populations that will be used.
  - `PopSize` determines the number of individuals inside each population.
  - `NGenerations` determines for many generations the populations will be evolved.
  - `HofSize` determines how many of the best individuals should be recorded.
  - `Model` is a struct that determines how to evolve each population of individuals.
- Optional fields
  - `ParallelEval` determines if a population is evaluated in parallel. The rule of thumb is to set this to `true` if your `Evaluate` method is expensive, if not it won't be worth the overhead. Refer to the [section on parallelism](#a-note-on-parallelism) for a more comprehensive explanation.
  - `Migrator` and `MigFrequency` should be provided if you want to exchange individuals between populations in case of a multi-population GA. If not the populations will be run independently. Again this is an advanced concept in the genetic algorithms field that you shouldn't deal with at first.
  - `Speciator` will split each population in distinct species at each generation. Each specie will be evolved separately from the others, after all the species has been evolved they are regrouped.
  - `Logger` can be used to record basic population statistics, you can read more about it in the [logging section](#logging-population-statistics).
  - `Callback` will execute any piece of code you wish every time `ga.Evolve()` is called. `Callback` will also be called when `ga.Initialize()` is. Using a callback can be useful for many things:
    - Calculating specific population statistics that are not provided by the logger
    - Changing parameters of the GA after a certain number of generations
    - Monitoring convergence
  - `EarlyStop` will be called before each generation to check if the evolution should be stopped early.
  - `RNG` can be set to make results reproducible. If it is not provided then a default `rand.New(rand.NewSource(time.Now().UnixNano()))` will be used. If you want to make your results reproducible use a constant source, e.g. `rand.New(rand.NewSource(42))`.

Once you have instantiated a `GAConfig` you can call it's `NewGA` method to obtain a `GA`. The `GA` struct has the following definition:

```go
type GA struct {
    GAConfig

    Populations Populations
    HallOfFame  Individuals
    Age         time.Duration
    Generations uint
}
```

Naturally a `GA` stores a copy of the `GAConfig` that was used to instantiate it. Apart from this the following fields are available:

- `Populations` is where all the current populations and individuals are kept.
- `HallOfFame` contains the `HofSize` best individuals ever encountered. This slice is always sorted, meaning that the first element of the slice will be the best individual ever encountered.
- `Age` indicates the duration the GA has spent evolving.
- `Generations` indicates how many how many generations have gone by.

You could bypass the `NewGA` method instantiate a `GA` with a `GAConfig` but this would leave the `GAConfig`'s fields unchecked for input errors.


#### Calling the Minimize method

You are now all set to find an optimal solution to your problem. To do so you have to call the GA's `Minimize` function which has the following signature:

```go
func (ga *GA) Minimize(newGenome func(rng *rand.Rand) Genome) error
```

You have to provide the `Minimize` a function which returns a `Genome`. It is recommended that the `Genome` thus produced contains random values. This is where the connection between the `Genome` interface and the `GA` struct is made.

The `Minimize` function will return an error (`nil` if everything went okay) once it is done. You can done access the first entry in the `HallOfFame` field to retrieve the best encountered solution.


#### Using the Slice interface

Classically GAs are used to optimize problems where the genome has a slice representation - eg. a vector or a sequence of DNA code. Almost all the mutation and crossover algorithms available in eaopt are based on the `Slice` interface which has the following definition.

```go
type Slice interface {
    At(i int) interface{}
    Set(i int, v interface{})
    Len() int
    Swap(i, j int)
    Slice(a, b int) Slice
    Split(k int) (Slice, Slice)
    Append(Slice) Slice
    Replace(Slice)
    Copy() Slice
}
```

Internally `IntSlice`, `Float64Slice` and `StringSlice` implement this interface so that you can use the available operators for most use cases. If however you wish to use the operators with slices of a different type you will have to implement the `Slice` interface. Although there are many methods to implement, they are all trivial (have a look at [`slice.go`](slice.go) and the [TSP example](https://github.com/MaxHalford/eaopt-examples/tree/master/tsp_grid).


#### Models

eaopt makes it easy to use different so called *models*. Simply put, a models defines how a GA evolves a population of individuals through a sequence of genetic operators. It does so without considering whatsoever the intrinsics of the underlying operators. In a nutshell, an evolution model attempts to mimic evolution in the real world. **It's extremely important to choose a good model because it is usually the highest influence on the performance of a GA**.

##### Generational model

The generational model is one the, if not the most, popular models. Simply put it generates *n* offsprings from a population of size *n* and replaces the population with the offsprings. The offsprings are generated by selecting 2 individuals from the population and applying a crossover method to the selected individuals until the *n* offsprings have been generated. The newly generated offsprings are then optionally mutated before replacing the original population. Crossover generates two new individuals, thus if the population size isn't an even number then the second individual from the last crossover (individual *n+1*) won't be included in the new population.

<div align="center">
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vQrkFXTHkak2GiRpDarsEIDHnsFWqXd9A98Cq2UUIR1keyMSU8NUE8af7_87KiQnmCKKBEb0IiQVsZM/pub?w=960&h=720" alt="generational" width="70%" />
</div>

##### Steady state model

The steady state model differs from the generational model in that the entire population isn't replaced between each generations. Instead of adding the children of the selected parents into the next generation, the 2 best individuals out of the two parents and two children are added back into the population so that the population size remains constant. However, one may also replace the parents with the children regardless of their fitness. This method has the advantage of not having to evaluate the newly generated offsprings. Whats more, crossover often generates individuals who are sub-par but who have a lot of potential; giving individuals generated from crossover a chance can be beneficial on the long run.

<div align="center">
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vTTk7b1QS67CZTr7-ksBMlk_cIDhm2YMZjemmrhXbLei5_VgvXCsINCLu8uia3ea6Ouj9I3V5HcZUwS/pub?w=962&h=499" alt="steady-state" width="70%" />
</div>

##### Select down to size model

The select down to size model uses two selection rounds. The first one is similar to the one used in the generational model. Parents are selected to generate new individuals through crossover. However, the offsprings are then merged with the original population and a second selection round occurs to determine which individuals will survive to the next generation. Formally *m* offsprings are generated from a population of *n*, the *n+m* individuals are then "selected down to size" so that there only remains *n* individuals.

<div align="center">
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vSyXQLPkWOOffKfnTRcdwrKvHTN9rWvdqGVT1fC6vcXGJAQPzxQVmauYLhSd2Xh74vQMhBEnhrSt1od/pub?w=969&h=946" alt="select-down-to-size" width="70%" />
</div>

##### Ring model

In the ring model, crossovers are applied to neighbours in a one-directional ring topology. Two by the two neighbours generate 2 offsprings. The best out of the 4 individuals (2 parents + 2 offsprings) replaces the first neighbour.

<div align="center">
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vTCsgqnEXj4KCn_C7IxHZXSw9XMP3RK_YeW5AoVKUSRHzq6CIFlp7fbBA-DK9mtFV330kROwrEsP6tj/pub?w=960&h=625" alt="ring" width="70%" />
</div>

##### Mutation only

It's possible to run a GA without crossover simply by mutating individuals. This can be done with the `ModMutationOnly` struct. At each generation each individual is mutated. `ModMutationOnly` has a `strict` field to determine if the mutant should replace the initial individual only if it's fitness is lower.

#### Speciation

Clusters, also called species in the literature, are a partitioning of individuals into smaller groups of similar individuals. Programmatically a cluster is a list of lists that each contain individuals. Individuals inside each species are supposed to be similar. The similarity depends on a metric, for example it could be based on the fitness of the individuals. In the literature, speciation is also called *speciation*.

The purpose of a partitioning individuals is to apply genetic operators to similar individuals. In biological terms this encourages "incest" and maintains isolated species. For example in nature animals usually breed with local mates and don't breed with different animal species.

Using speciation/speciation with genetic algorithms became "popular" when they were first applied to the [optimization of neural network topologies](https://www.wikiwand.com/en/Neuroevolution_of_augmenting_topologies). By mixing two neural networks during crossover, the resulting neural networks were often useless because the inherited weights were not optimized for the new topology. This meant that newly generated neural networks were not performing well and would likely disappear during the selection phase. Thus speciation was introduced so that neural networks evolved in similar groups in order for new neural networks wouldn't disappear immediately. Instead the similar neural networks would evolve between each other until they were good enough to mixed with the other neural networks.

With eaopt it's possible to use speciation on top of all the rest. To do so the `Speciator` field of the `GA` struct has to specified.

<div align="center">
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vRLr7j4ML-ZeXFfvjko9aepRAkCgBlpg4dhuWhB-vXCQ17gJFmDQHrcUbcPFwlqzvaPAXwDxx5ld1kf/pub?w=686&h=645" alt="speciation" width="70%" />
</div>

#### Multiple populations and migration

Multi-populations GAs run independent populations in parallel. They are not frequently used, however they are very easy to understand and to implement. In eaopt a `GA` struct contains a `Populations` field which stores each population in a slice. The number of populations is specified in the `GAConfig`'s `NPops` field.

If `Migrator` and `MigFrequency` are not provided the populations will be run independently in parallel. However, if they are provided then at each generation number that is divisible by `MigFrequency` (for example 5 divides generation number 25) individuals will be exchanged between the populations following the `Migrator`.

Using multi-populations can be an easy way to gain in diversity. Moreover, not using multi-populations on a multi-core architecture is a waste of resources.

With eaopt you can use multi-populations and speciation at the same time. The following flowchart shows what that would look like.

<div align="center">
  <img src="https://docs.google.com/drawings/d/14VVpTkWquhrcG_oQ61hvZgjKlYWZs_UZRVnL22HFYKM/pub?w=1052&h=607" alt="multi-population_and_speciation" width="70%" />
</div>


#### Logging population statistics

It's possible to log statistics for each population at every generation. To do so you simply have to provide the `GA` struct a `Logger` from the Go standard library. This is quite convenient because it allows you to decide where to write the log output, whether it be in a file or directly in the standard output.

```go
ga.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
```

If a logger is provided, each row in the log output will include

- the population ID,
- the population minimum fitness,
- the population maximum fitness,
- the population average fitness,
- the population's fitness standard deviation.

### Particle swarm optimization

#### Description

[Particle swarm optimization (PSO)](https://www.wikiwand.com/en/Particle_swarm_optimization) can be used to optimize real valued functions. It maintains a population of candidate solutions called particles. The particles move around the search-space according to a mathematical formula that takes as input the particle's position and it's velocity. Each particle's movement is influenced by its's local best encountered position, as well as the best overall position in the search-space (these values are updated after each generation). This is expected to move the swarm toward the best solutions.

As can be expected there are many variants of PSO. The `SPSO` struct implements the [SPSO-2011 standard](http://clerc.maurice.free.fr/pso/SPSO_descriptions.pdf).

#### Example

In this example we're going to minimize th [Styblinski-Tang function](https://www.sfu.ca/~ssurjano/stybtang.html) with two dimensions. The global minimum is about -39.16599 times the number of dimensions.

```go
package main

import (
    "fmt"
    m "math"
    "math/rand"

    "github.com/MaxHalford/eaopt"
)

func StyblinskiTang(X []float64) (y float64) {
    for _, x := range X {
        y += m.Pow(x, 4) - 16*m.Pow(x, 2) + 5*x
    }
    return 0.5 * y
}

func main() {
    // Instantiate SPSO
    var spso, err = eaopt.NewDefaultSPSO()
    if err != nil {
        fmt.Println(err)
        return
    }

    // Fix random number generation
    spso.GA.RNG = rand.New(rand.NewSource(42))

    // Run minimization
    _, y, err := spso.Minimize(StyblinskiTang, 2)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Output best encountered solution
    fmt.Printf("Found minimum of %.5f, the global minimum is %.5f\n", y, -39.16599*2)
}
```

This should produce the following output.

```sh
>>> Found minimum of -78.23783, the global minimum is -78.33198
```

#### Parameters

You can (and should) instantiate an `SPSO` with the `NewSPSO` method. You can also use the `NewDefaultSPSO` method as is done in the previous example.

```go
func NewSPSO(nParticles, nSteps uint, min, max, w float64, parallel bool, rng *rand.Rand) (*SPSO, error)
```

- `nParticles` is the number of particles to use
- `nSteps` is the number of steps during which evolution occurs
- `min` and `max` are the boundaries from which the initial values are sampled from
- `w` is the velocity amplifier
- `parallel` determines if the particles are evaluated in parallel or not
- `rng` is a random number generator, you can set it to `nil` if you want it to be random

### Differential evolution

#### Description

[Differential evolution (DE)](https://www.wikiwand.com/en/Differential_evolution) somewhat resembles PSO and is also used for optimizing real-valued functions. At each generation, each so-called agent is moved according to the position of 3 randomly sampled agents. If the new position is not better than the current one then it is discarded.

As can be expected there are many variants of PSO. The `SPSO` struct implements the [SPSO-2011 standard](http://clerc.maurice.free.fr/pso/SPSO_descriptions.pdf).

#### Example

In this example we're going to minimize th [Ackley function](https://www.sfu.ca/~ssurjano/ackley.html) with two dimensions. The global minimum is 0.

```go
package main

import (
    "fmt"
    m "math"
    "math/rand"

    "github.com/MaxHalford/eaopt"
)

func Ackley(x []float64) float64 {
    var (
        a, b, c = 20.0, 0.2, 2 * m.Pi
        s1, s2  float64
        d       = float64(len(x))
    )
    for _, xi := range x {
        s1 += xi * xi
        s2 += m.Cos(c * xi)
    }
    return -a*m.Exp(-b*m.Sqrt(s1/d)) - m.Exp(s2/d) + a + m.Exp(1)
}

func main() {
    // Instantiate DiffEvo
    var de, err = eaopt.NewDefaultDiffEvo()
    if err != nil {
        fmt.Println(err)
        return
    }

    // Fix random number generation
    de.GA.RNG = rand.New(rand.NewSource(42))

    // Run minimization
    _, y, err := de.Minimize(Ackley, 2)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Output best encountered solution
    fmt.Printf("Found minimum of %.5f, the global minimum is 0\n", y)
}
```

This should produce the following output.

```sh
>>> Found minimum of 0.00137, the global minimum is 0
```

#### Parameters

You can (and should) instantiate an `DiffEvo` with the `NewDiffEvo` method. You can also use the `NewDefaultDiffEvo` method as is done in the previous example.

```go
func NewDiffEvo(nAgents, nSteps uint, min, max, cRate, dWeight float64, parallel bool, rng *rand.Rand) (*DiffEvo, error)
```

- `nAgents` is the number of agents to use (it has to be at least 4)
- `nSteps` is the number of steps during which evolution occurs
- `min` and `max` are the boundaries from which the initial values are sampled from
- `cRate` is the crossover rate
- `dWeight` is the differential weight
- `parallel` determines if the agents are evaluated in parallel or not
- `rng` is a random number generator, you can set it to `nil` if you want it to be random


### OpenAI evolution strategy

#### Description

[OpenAI](https://openai.com/) proposed [a simple evolution strategy](https://blog.openai.com/evolution-strategies/) based on the use of natural gradients. The algorithm is dead simple:

1. Choose a center `mu` at random
2. Sample points around `mu` using a normal distribution
3. Evaluate each point and obtain the natural gradient `g`
4. Move `mu` along the natural gradient `g` using a learning rate
5. Repeat from step 2 until satisfied

#### Example

In this example we're going to minimize th [Rastrigin function](https://www.sfu.ca/~ssurjano/rastr.html) with three dimensions. The global minimum is 0.

```go
package main

import (
    "fmt"
    m "math"
    "math/rand"

    "github.com/MaxHalford/eaopt"
)

func Rastrigin(x []float64) (y float64) {
    y = 10 * float64(len(x))
    for _, xi := range x {
        y += m.Pow(xi, 2) - 10*m.Cos(2*m.Pi*xi)
    }
    return y
}

func main() {
    // Instantiate DiffEvo
    var oes, err = eaopt.NewDefaultOES()
    if err != nil {
        fmt.Println(err)
        return
    }

    // Fix random number generation
    oes.GA.RNG = rand.New(rand.NewSource(42))

    // Run minimization
    _, y, err := oes.Minimize(Rastrigin, 2)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Output best encountered solution
    fmt.Printf("Found minimum of %.5f, the global minimum is 0\n", y)
}
```

This should produce the following output.

```sh
>>> Found minimum of 0.02270, the global minimum is 0
```

#### Parameters

You can (and should) instantiate an `OES` with the `NewOES` method. You can also use the `NewDefaultOES` method as is done in the previous example.

```go
func NewOES(nPoints, nSteps uint, sigma, lr float64, parallel bool, rng *rand.Rand) (*OES, error)
```

- `nPoints` is the number of points to use (it has to be at least 3)
- `nSteps` is the number of steps during which evolution occurs
- `sigma` determines the shape of the normal distribution used to sample new points
- `lr` is the learning rate
- `parallel` determines if the agents are evaluated in parallel or not
- `rng` is a random number generator, you can set it to `nil` if you want it to be random



## A note on parallelism

Evolutionary algorithms are famous for being [embarrassingly parallel](https://www.wikiwand.com/en/Embarrassingly_parallel). Most of the operations can be run independently each one from another. For example individuals can be mutated in parallel because mutation doesn't have any side effects.

The Go language provides nice mechanisms to run stuff in parallel, provided you have more than one core available. However, parallelism is only worth it when the functions you want to run in parallel are heavy. If the functions are cheap then the overhead of spawning routines will be too high and not worth it. It's simply not worth using a routine for each individual because operations at an individual level are often not time consuming enough.

By default eaopt will evolve populations in parallel. This is because evolving one population implies a lot of operations and parallelism is worth it. If your `Evaluate` method is heavy then it might be worth evaluating individuals in parallel, which can done by setting the `GA`'s `ParallelEval` field to `true`. Evaluating individuals in parallel can be done regardless of the fact that you are using more than one population.


## FAQ

**What if I don't want to use crossover?**

Alas you still have to implement the `Genome` interface. You can however provide a blank `Crossover` method just to satisfy the interface.

```go
type Vector []float64

func (X Vector) Crossover(Y eaopt.Genome, rng *rand.Rand) {}
```


**Why aren't my `Mutate` and `Crossover` methods modifying my `Genome`s?**

The `Mutate` and `Crossover` methods have to modify the values of the `Genome` in-place. The following code will work because the `Vector` is a slice; slices in Go are references to underlying data, hence modifying a slice modifies them in-place.

```go
type Vector []float64

func (X Vector) Mutate(rng *rand.Rand) {
    eaopt.MutNormal(X, rng, 0.5)
}
```

On the contrary, mutating other kind of structs will require the `*` symbol to access the struct's pointer. Notice the `*Name` in the following example.

```go
type Name string

func (n *Name) Mutate(rng *rand.Rand) {
    n = randomName()
}
```


**Are evolutionary optimization algorithms any good?**

For real-valued, differentiable functions, evolutionary optimization algorithms will probably not fair well against methods  based on gradient descent. Intuitively this is because evolutionary optimization algorithms ignore the shape and slope of the function. However gradient descent algorithms usually get stuck in local optimas, whereas evolutionary optimization algorithms don't.

As mentioned earlier, some problems can simply not be written down as [closed-form expressions](https://www.wikiwand.com/en/Closed-form_expression). In this case methods based on gradient information can't be used whilst evolutionary optimization algorithms can still be used. For example tuning the number of layers and of neurons per layer in a neural network is an open problem that doesn't yet have a reliable solution. Neural networks architectures used in production are usually designed by human experts. The field of [neuroevolution](https://www.wikiwand.com/en/Neuroevolution) aims to train neural networks with evolutionary algorithms.

**How can I contribute?**

Feel free to implement your own operators or to make suggestions! Check out the [CONTRIBUTING file](CONTRIBUTING.md) for some guidelines. [This repository](https://github.com/fcampelo/EC-Bestiary) has a long list of existing evolutionary algorithms.


## Dependencies

You can see the list of dependencies [here](https://godoc.org/github.com/MaxHalford/eaopt?imports) and the graph view [here](https://godoc.org/github.com/MaxHalford/eaopt?import-graph). Here is the list of external dependencies:

- [golang.org/x/sync/errgroup](https://godoc.org/golang.org/x/sync/errgroup)


## License

The MIT License (MIT). Please see the [LICENSE file](LICENSE) for more information.
