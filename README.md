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
  <a href="https://godoc.org/github.com/MaxHalford/gago">
    <img src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square" alt="godoc" />
  </a>
  <!-- Build status -->
  <a href="https://travis-ci.org/MaxHalford/gago">
    <img src="https://img.shields.io/travis/MaxHalford/gago/master.svg?style=flat-square" alt="build_status" />
  </a>
  <!-- Test coverage -->
  <a href="https://coveralls.io/github/MaxHalford/gago?branch=master">
    <img src="https://coveralls.io/repos/github/MaxHalford/gago/badge.svg?branch=master&style=flat-square" alt="test_coverage" />
  </a>
  <!-- Go report card -->
  <a href="https://goreportcard.com/report/github.com/MaxHalford/gago">
    <img src="https://goreportcard.com/badge/github.com/MaxHalford/gago?style=flat-square" alt="go_report_card" />
  </a>
  <!-- Code Climate -->
  <a href="https://codeclimate.com/github/MaxHalford/gago">
    <img src="https://codeclimate.com/github/MaxHalford/gago/badges/gpa.svg" alt="Code Climate" />
  </a>
  <!-- License -->
  <a href="https://opensource.org/licenses/MIT">
    <img src="http://img.shields.io/:license-mit-ff69b4.svg?style=flat-square" alt="license"/>
  </a>
</div>

<br/>

<div align="center">An extensible toolkit for conceiving and running genetic algorithms</div>

<br/>

**Table of Contents**
  - [Example](#example)
  - [Background](#background)
    - [Terminology](#terminology)
    - [Methodology](#methodology)
  - [Features](#features)
  - [Usage](#usage)
    - [Implementing the Genome interface](#implementing-the-genome-interface)
    - [Using the Slice interface](#using-the-slice-interface)
    - [Instantiating a GA struct](#instantiating-a-ga-struct)
    - [Running a GA](#running-a-ga)
    - [Models](#models)
    - [Multiple populations and migration](#multiple-populations-and-migration)
    - [Speciation](#speciation)
    - [Presets](#presets)
    - [Logging population statistics](#logging-population-statistics)
  - [A note on parallelism](#a-note-on-parallelism)
  - [FAQ](#faq)
  - [Alternatives](#alternatives)
  - [License](#license)

## Example

The following example attempts to minimize the [Drop-Wave function](https://www.sfu.ca/~ssurjano/drop.html) which is known to have a minimum value of -1 when each of it's arguments is equal to 0.

<div align="center">
  <img src="https://github.com/MaxHalford/gago-examples/blob/master/drop_wave/chart.png" alt="drop_wave_chart" />
  <img src="https://github.com/MaxHalford/gago-examples/blob/master/drop_wave/function.png" alt="drop_wave_function" />
</div>

```go
package main

import (
    "fmt"
    m "math"
    "math/rand"

    "github.com/MaxHalford/gago"
)

// A Vector contains float64s.
type Vector []float64

// Evaluate a Vector with the Drop-Wave function which takes two variables as
// input and reaches a minimum of -1 in (0, 0). The function is rather pure so
// there isn't any error handling to do.
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
    gago.MutNormalFloat64(X, 0.8, rng)
}

// Crossover a Vector with another Vector by applying uniform crossover.
func (X Vector) Crossover(Y gago.Genome, rng *rand.Rand) {
    gago.CrossUniformFloat64(X, Y.(Vector), rng)
}

// Clone a Vector to produce a new one that points to a different slice.
func (X Vector) Clone() gago.Genome {
    var Y = make(Vector, len(X))
    copy(Y, X)
    return Y
}

// VectorFactory returns a random vector by generating 2 values uniformally
// distributed between -10 and 10.
func VectorFactory(rng *rand.Rand) gago.Genome {
    return Vector(gago.InitUnifFloat64(2, -10, 10, rng))
}

func main() {
    var ga = gago.Generational(VectorFactory)
    var err = ga.Initialize()
    if err != nil {
        fmt.Println("Handle error!")
    }

    fmt.Printf("Best fitness at generation 0: %f\n", ga.HallOfFame[0].Fitness)
    for i := 1; i < 10; i++ {
        err = ga.Evolve()
        if err != nil {
            fmt.Println("Handle error!")
        }
        fmt.Printf("Best fitness at generation %d: %f\n", i, ga.HallOfFame[0].Fitness)
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
```

**More examples**

All the examples can be found [here](https://github.com/MaxHalford/gago-examples).

- [Cross-in-Tray (speciation)](https://github.com/MaxHalford/gago-examples/tree/master/cross_in_tray)
- [Grid TSP](https://github.com/MaxHalford/gago-examples/tree/master/tsp_grid)
- [Including a constraint](https://github.com/MaxHalford/gago-examples/tree/master/constraint)
- [One Max problem](https://github.com/MaxHalford/gago-examples/tree/master/one_max)
- [N-queens problem](https://github.com/MaxHalford/gago-examples/tree/master/n_queens)
- [String matching](https://github.com/MaxHalford/gago-examples/tree/master/string_matching)

## Background

There is a lot of intellectual fog around the concept of genetic algorithms (GAs). It's important to appreciate the fact that GAs are composed of many nuts and bolts. **There isn't a single definition of genetic algorithms**. `gago` is intended to be a toolkit where one may run many kinds of genetic algorithms, with different evolution models and various genetic operators.

### Terminology

- ***Fitness function***: The fitness function is simply the function associated to a given problem. It takes in an input and returns an output. The goal is to find the input that minimizes the output.
- ***Individual***: An individual contains a **genome** which represents a function input. In the physical world, an individual's genome is composed of acids. In an imaginary world, it could be composed of floating point numbers or string sequences representing cities. A **fitness** can be associated to a genome thanks to the fitness function. For example, one could measure the height of a group of individuals in order to rank them. In this case the genome is the body of an individual, the fitness function is the act of measuring the height of an individual's body and the fitness is the resulting height.
- ***Population***: Individuals are contained in a population wherein they can interact.
- ***Crossover***:  A crossover acts on two or more individuals (called **parents**) and mixes their genome in order to produce one or more new individuals (called **offsprings**). Crossover is really what sets genetic algorithms apart from other evolutionary methods.
- ***Selection***: Selection is a process in which parents are selected to generate offsprings, most often by applying a crossover method. Popular selection methods include **elitism selection** and **tournament selection**.
- ***Mutation***: Mutation applies random modifications to an individual's genome without interacting with other individuals.
- ***Migration***: **Multi-population** GAs run more than one population in parallel and exchange individuals between each other.
- ***Speciation***: In the physical world, individuals do not mate at random. Instead, they mate with similar individuals. For some problems -- for example neural network topology optimization -- crossover will often generate poor solutions. Speciation sidesteps this by mating similar individuals (called **species**) separately.
- ***Evolution model***: An evolution model describes the exact manner and order in which genetic operators are applied to a population. The most popular models are the **steady state model** and the **generational model**.

### Methodology

In a nutshell, a GA solves an optimization problem by doing the following:

1. Generate random solutions.
2. Assign a fitness to each solutions.
3. Sort the solutions according to their fitness.
4. Apply genetic operators following a model.
5. Repeat from step 2 until the stopping criterion is satisfied.

This description is voluntarily vague as to how the genetic operators are applied. It's important to understand that there isn't a single way of applying genetic algorithms. For example some people believe that crossover is useless and use mutation for generating new individuals. In general genetic operators are applied following an **evolution model**, a fact that is often omitted in introductions to genetic algorithms. Popular stopping criteria include

- a fixed number of generations,
- a fixed duration,
- an indicator that the population is stagnating.


## Features

- gago is extensible, you can control most of the evolution logic
- Different evolution models are available
- Popular operators are already implemented
- Speciation is available
- Multiple population migration is available


## Usage

The two requirements for using gago are

- Implement the `Genome` interface.
- Instantiate a `GA` struct.

The `Genome` interface is used to define the logic that is specific to your problem; logic that gago doesn't know about. For example this is where you will define an `Evaluate()` method for evaluating a particular problem. The `GA` struct contains context-agnostic information. For example this is where you can choose the number of individuals in a population (which is a separate concern from your particular problem). Apart from a good design pattern, decoupling the problem definition from the optimization through the `Genome` interface means that gago can be used to optimize *any* kind of problem.

### Implementing the Genome interface

Let's have a look at the `Genome` interface.

```go
type Genome interface {
    Evaluate() (float64, error)
    Mutate(rng *rand.Rand)
    Crossover(genome Genome, rng *rand.Rand)
    Clone() Genome
}
```

The `Evaluate()` method returns the fitness of a genome. The sweet thing is that you can do whatever you want in this method. Your struct that implements the interface doesn't necessarily have to be a slice. The `Evaluate()` method is *your* problem to deal with. gago only needs it's output to be able to function. You can also return an `error` which gago will catch and return when calling `ga.Initialize()` and `ga.Evolve()`.

The `Mutate(rng *rand.Rand)` method is where you can modify an existing solution by tinkering with it's variables. The way in which you should mutate a solution essentially boils down to your particular problem. gago provides some common mutation methods that you can use instead of reinventing the wheel -- this is what is being done in most of the [examples](https://github.com/MaxHalford/gago-examples).

The `Crossover(genome Genome, rng *rand.Rand)` method combines two individuals. The important thing to notice is that the type of first argument differs from the struct calling the method. Indeed the first argument is a `Genome` that has to be casted into your struct before being able to apply a crossover operator. This is due to the fact that Go doesn't provide generics out of the box; it's easier to convince yourself by checking out the examples.

The `Clone()` method is there to produce independent copies of the struct you want to evolve. This is necessary for internal reasons and ensures that pointer fields are not pointing to identical memory addresses. Usually this is not too difficult implement; you just have to make sure that the clones you produce are not shallow copies of the genome that is being cloned. This is also fairly easy to unit test.

Once you have implemented the `Genome` interface you have provided gago with all the information it couldn't guess for you. Essentially you have total control over the definition of your problem, gago will handle the rest and find a good solution to the problem.

### Using the Slice interface

Classically GAs are used to optimize problems where the genome has a slice representation - eg. a vector or a sequence of DNA code. Almost all the mutation and crossover algorithms available in gago are based on the `Slice` interface which has the following definition.

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

Internally `IntSlice`, `Float64Slice` and `StringSlice` implement this interface so that you can use the available operators for most use cases. If however you wish to use the operators which slices of a different type you will have to implement the `Slice` interface. Although there are many methods to implement, they are all trivial (have a look at [`slice.go`](slice.go) and the [TSP example](https://github.com/MaxHalford/gago-examples/tree/master/tsp_grid).


### Instantiating a GA struct

Let's have a look at the GA struct.

```go
type GA struct {
    // Required fields
    NewGenome    NewGenome
    NPops        int
    PopSize      int
    Model        Model

    // Optional fields
    NBest        int
    Migrator     Migrator
    MigFrequency int
    Speciator    Speciator
    Logger       *log.Logger
    Callback     func(ga *GA)
    RNG          *rand.Rand
    ParallelEval bool

    // Fields generated at runtime
    Populations        Populations
    HallOfFame         Individuals
    Age                time.Duration
    Generations        int
}
```

You have to fill in the first set of fields, the rest are generated when calling the `GA`'s `Initialize()` method. Check out the examples in `presets.go` to get an idea of how to fill them out.

- Required fields
  - `NewGenome` is a method that returns a random `Genome` that you defined in the previous step. gago will use this method to produce initial individuals in each population. Again, gago provides some methods for common random genome generation.
  - `NPops` determines the number of populations that will be used.
  - `PopSize` determines the number of individuals inside each population.
  - `Model` determines how to use the genetic operators you chose in order to produce better solutions, in other words it's a recipe. A dedicated section is available in the [model section](#models).
- Optional fields
  - `NBest` determines how many of the best individuals encountered should be regarded in the `HallOfFame` field. This defaults to 1.
  - `Migrator` and `MigFrequency` should be provided if you want to exchange individuals between populations in case of a multi-population GA. If not the populations will be run independently. Again this is an advanced concept in the genetic algorithms field that you shouldn't deal with at first.
  - `Speciator` will split each population in distinct species at each generation. Each specie will be evolved separately from the others, after all the species has been evolved they are regrouped.
  - `Logger` is optional and provides basic population statistics, you can read more about it in the [logging section](#logging-population-statistics).
  - `Callback` is optional will execute any piece of code you wish every time `ga.Evolve()` is called. `Callback` will also be called when `ga.Initialize()` is. Using a callback can be useful for many things:
    - Calculating specific population statistics that are not provided by the logger
    - Changing parameters of the GA after a certain number of generations
    - Monitoring for converging populations
  - `RNG` can be set to make results reproducible. If it is not provided then a default `rand.New(rand.NewSource(time.Now().UnixNano()))` will be used. If you want to make your results reproducible use a constant source, e.g. `rand.New(rand.NewSource(42))`.
  - `ParallelEval` determines if a population is evaluated in parallel. The rule of thumb is to set this to `true` if your `Evaluate` method is expensive, if not it won't be worth the overhead. Refer to the [section on parallelism](#a-note-on-parallelism) for a more comprehensive explanation.
- Fields populated at runtime
  - `Populations` is where all the current populations and individuals are kept.
  - `HallOfFame` contains the `NBest` individuals ever encountered. This slice is always sorted, meaning that the first element of the slice will be the best individual ever encountered.
  - `Age` indicates the duration the GA has spent calling the `Evolve` method.
  - `Generations` indicates how many times the `Evolve` method has been called.


### Running a GA

Once you have implemented the `Genome` interface and instantiated a `GA` struct you are good to go. You can call the `GA`'s `Evolve` method which will apply a model once (see the [models section](#models)). It's your choice if you want to call `Evolve` method multiple by using a loop or by imposing a time limit. The `Evolve` method will return an `error` which you should handle. If your population does not evolve when you call `Evolve` it's most likely because `Evolve` returned an error.

At any time you have access to the `GA`'s `HallOfFame` field which contains the `NBest` individuals ever encountered.

### Models

`gago` makes it easy to use different so called *models*. Simply put, a models tells the story of how a GA enhances a population of individuals through a sequence of genetic operators. It does so without considering whatsoever the underlying operators. In a nutshell, an evolution model attempts to mimic evolution in the real world. **It's extremely important to choose a good model because it is usually the highest influence on the performance of a GA**.

#### Generational model

The generational model is one the, if not the most, popular models. Simply put it generates *n* offsprings from a population of size *n* and replaces the population with the offsprings. The offsprings are generated by selecting 2 individuals from the population and applying a crossover method to the selected individuals until the *n* offsprings have been generated. The newly generated offsprings are then optionally mutated before replacing the original population. Crossover generates two new individuals, thus if the population size isn't an even number then the second individual from the last crossover (individual *n+1*) won't be included in the new population.

<div align="center">
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vQrkFXTHkak2GiRpDarsEIDHnsFWqXd9A98Cq2UUIR1keyMSU8NUE8af7_87KiQnmCKKBEb0IiQVsZM/pub?w=960&h=720" alt="generational" width="70%" />
</div>

#### Steady state model

The steady state model differs from the generational model in that the entire population isn't replaced between each generations. Instead of adding the children of the selected parents into the next generation, the 2 best individuals out of the two parents and two children are added back into the population so that the population size remains constant. However, one may also replace the parents with the children regardless of their fitness. This method has the advantage of not having to evaluate the newly generated offsprings. Whats more, crossover often generates individuals who are sub-par but who have a lot of potential; giving individuals generated from crossover a chance can be beneficial on the long run.

<div align="center">
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vTTk7b1QS67CZTr7-ksBMlk_cIDhm2YMZjemmrhXbLei5_VgvXCsINCLu8uia3ea6Ouj9I3V5HcZUwS/pub?w=962&h=499" alt="steady-state" width="70%" />
</div>

#### Select down to size model

The select down to size model uses two selection rounds. The first one is similar to the one used in the generational model. Parents are selected to generate new individuals through crossover. However, the offsprings are then merged with the original population and a second selection round occurs to determine which individuals will survive to the next generation. Formally *m* offsprings are generated from a population of *n*, the *n+m* individuals are then "selected down to size" so that there only remains *n* individuals.

<div align="center">
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vSyXQLPkWOOffKfnTRcdwrKvHTN9rWvdqGVT1fC6vcXGJAQPzxQVmauYLhSd2Xh74vQMhBEnhrSt1od/pub?w=969&h=946" alt="select-down-to-size" width="70%" />
</div>

#### Ring model

In the ring model, crossovers are applied to neighbours in a one-directional ring topology. Two by the two neighbours generate 2 offsprings. The best out of the 4 individuals (2 parents + 2 offsprings) replaces the first neighbour.

<div align="center">
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vTCsgqnEXj4KCn_C7IxHZXSw9XMP3RK_YeW5AoVKUSRHzq6CIFlp7fbBA-DK9mtFV330kROwrEsP6tj/pub?w=960&h=625" alt="ring" width="70%" />
</div>

#### Simulated annealing

Although [simulated annealing](https://www.wikiwand.com/en/Simulated_annealing) isn't a genetic algorithm, it can nonetheless be implemented with gago. A mutator is the only necessary operator. Other than that a starting temperature, a stopping temperature and a decrease rate have to be provided. Effectively a single simulated annealing is run for each individual in the population.

The temperature evolution is relative to one single generation. In order to mimic the original simulated annealing algorithm, one would the number of individuals to 1 and would run the algorithm for only 1 generation. However, nothing stops you from running many simulated annealings and to repeat them over many generations.

#### Mutation only

It's possible to run a GA without crossover simply by mutating individuals. Essentially this boils down to doing [hill climbing](https://www.wikiwand.com/en/Hill_climbing) because there is not interaction between individuals. Indeed taking a step in hill climbing is equivalent to mutation for genetic algorithms. What's nice is that by using a population of size n you are essentially running multiple independent hill climbs.

### Speciation

Clusters, also called speciation in the literature, are a partitioning of individuals into smaller groups of similar individuals. Programmatically a cluster is a list of lists each containing individuals. Individuals inside each species are supposed to be similar. The similarity depends on a metric, for example it could be based on the fitness of the individuals. In the literature, speciation is also called *speciation*.

The purpose of a partitioning individuals is to apply genetic operators to similar individuals. In biological terms this encourages "incest" and maintains isolated species. For example in nature animals usually breed with local mates and don't breed with different animal species.

Using speciation/speciation with genetic algorithms became "popular" when they were first applied to the [optimization of neural network topologies](https://www.wikiwand.com/en/Neuroevolution_of_augmenting_topologies). By mixing two neural networks during crossover, the resulting neural networks were often useless because the inherited weights were not optimized for the new topology. This meant that newly generated neural networks were not performing well and would likely disappear during selection. Thus speciation was introduced so that neural networks evolved in similar groups so that new neural networks wouldn't disappear immediately. Instead the similar neural networks would evolve between each other until they were good enough to mixed with the other neural networks.

With gago it's possible to use speciation on top of all the rest. To do so the `Speciator` field of the `GA` struct has to specified.

<div align="center">
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vRLr7j4ML-ZeXFfvjko9aepRAkCgBlpg4dhuWhB-vXCQ17gJFmDQHrcUbcPFwlqzvaPAXwDxx5ld1kf/pub?w=686&h=645" alt="speciation" width="70%" />
</div>

### Multiple populations and migration

Multi-populations GAs run independent populations in parallel. They are not frequently used, however they are very easy to understand and to implement. In gago a `GA` struct contains a `Populations` field which contains each population. The number of populations is specified in the `GA`'s `NPops` field.

If `Migrator` and `MigFrequency` are not provided the populations will be run independently, in parallel. However, if they are provided then at each generation number that divides `MigFrequency` (for example 5 divides 25) individuals will be exchanged between the populations following the `Migrator` protocol.

Using multi-populations can be an easy way to gain in diversity. Moreover, not using multi-populations on a multi-core architecture is a waste of resources.

With gago you can use multi-populations and speciation at the same time. The following flowchart shows what that would look like.

<div align="center">
  <img src="https://docs.google.com/drawings/d/14VVpTkWquhrcG_oQ61hvZgjKlYWZs_UZRVnL22HFYKM/pub?w=1052&h=607" alt="multi-population_and_speciation" width="70%" />
</div>

### Presets

Some preset GA instances are available to get started as fast as possible. They are available in the [presets.go](presets.go) file. These instances also serve as example instantiations of the GA struct. To obtain optimal solutions you should fill in the fields manually!

### Logging population statistics

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


## A note on parallelism

Genetic algorithms are famous for being [embarrassingly parallel](https://www.wikiwand.com/en/Embarrassingly_parallel). Most of the operations used in the GA can be run independently each one from another. For example individuals can be mutated in parallel because mutation doesn't have any side effects.

The Go language provides nice mechanisms to run stuff in parallel, provided you have more than one core available. However, parallelism is only worth it when the functions you want to run in parallel are "heavy". If the functions are cheap then the overhead of spawning routines will be too high and not worth it. It's simply not worth using a routine for each individual because operations at an individual level are often not time consuming enough.

By default gago will evolve populations in parallel. This is because evolving one population implies a lot of operations and parallelism is worth it. If your `Evaluate` method is heavy then it might be worth evaluating individuals in parallel, which can done by setting the `GA`'s `ParallelEval` field to `true`. Evaluating individuals in parallel can be done regardless of the fact that you are using more than one population.


## FAQ

**What if I don't want to use crossover?**

Alas you still have to implement the `Genome` interface. You can however provide a blank `Crossover` method just to satisfy the interface.

```go
type Vector []float64

func (X Vector) Crossover(Y gago.Genome, rng *rand.Rand) (gago.Genome, gago.Genome) {
}
```

Or you can set `CrossRate: 0` when initializing the GA model.


**Why aren't my `Mutate` and `Crossover` methods modifying my `Genome`s?**

The `Mutate` and `Crossover` methods have to modify the values of the `Genome` in-place. The following code will work because the `Vector` is a slice; slices in Go are references to underlying data, hence modifying a slice modifies them in-place.

```go
type Vector []float64

func (X Vector) Mutate(rng *rand.Rand) {
    gago.MutNormal(X, rng, 0.5)
}
```

On the contrary, mutating other kind of structs will require the `*` symbol to access the struct's pointer. Notice the `*Name` in the following example.

```go
type Name string

func (n *Name) Mutate(rng *rand.Rand) {
    n = randomName()
}
```


**Should I be using genetic algorithms?**

Genetic algorithms (GAs) are often used for [NP-hard problems](https://www.wikiwand.com/en/NP-hardness). They *usually* perform better than [hill climbing](https://www.wikiwand.com/en/Hill_climbing) and [simulated annealing](https://www.wikiwand.com/en/Simulated_annealing) because they explore the search space more intelligently. However, GAs can also be used for classical problems where the search space makes it difficult for, say, gradient algorithms to be efficient (like the introductory example).

As mentioned earlier, some problems can simply not be written down as [closed-form expressions](https://www.wikiwand.com/en/Closed-form_expression). For example tuning the number of layers and of neurons per layer in a neural network is an open problem that doesn't yet have a reliable solution. Neural networks architectures used in production are usually designed by human experts. The field of [neuroevolution](https://www.wikiwand.com/en/Neuroevolution) aims to train neural networks with evolutionary algorithms. As such genetic algorithms are a good candidate for training neural networks, usually by optimizing the network's topology.

**How can I contribute?**

Feel free to implement your own operators or to make suggestions! Check out the [CONTRIBUTING file](CONTRIBUTING.md) for some guidelines.


## Dependencies

You can see the list of dependencies [here](https://godoc.org/github.com/MaxHalford/gago?imports) and the graph view [here](https://godoc.org/github.com/MaxHalford/gago?import-graph). Here is the list of external dependencies:

- [golang.org/x/sync/errgroup](https://godoc.org/golang.org/x/sync/errgroup)

## License

The MIT License (MIT). Please see the [LICENSE file](LICENSE) for more information.
