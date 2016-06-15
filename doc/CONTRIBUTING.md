# Contributing to `gago`

## Genetic operators

There are many genetic operators yet to be implemented. Feel free to send pull requests with your implementations. The only requirements are that the genetic operator respects the existing naming conventions and includes a test in the corresponding `*_test.go` file.

## Roadmap

- Settings verification
- Tidy tests
- Statistics
- Benchmarking vs other libraries
- Profiling
- Implement more genetic operators
- More examples
- Possibility to add manual heuristics
- Add different stopping criterions
- Count evaluations

## Code style

- Keep names short
- Use `var x = "hello"` instead of `x := "hello"`
- Use a single `var` for consecutive variable assignments.
- Aim for 80 characters per line

## Naming conventions

### Genetic operators

There is a convention for naming genetic operators. The name begins with an abreviation of the kind of operator being implemented:

- `Cross`: crossover
- `Init`: initializer
- `Mig`: migrator
- `Mut`: mutator
- `Sel`: selector
- `Clu`: clusterer

Then comes the second part of the name which indicates the name of the genetic operator.

Finally the name of the operator ends with a letter indicating on what kind of genomes the operator works:

- `F`: `float64`
- `S`: `string`
- No letter means the operator works on any kind of genome, regardless of the
underlying type.

For example `MutUniformS` is a *mutator* operator, it applies *uniform* mutation
on a *string* genome.

### Shortnames

Along with the genetic operators prefixes, other shortnames are used:

- `indi`: individual
- `pop`: population
- `ff`: fitness function

### Models

Implementations of the `Model` interface begin with a `Mod` prefix.


## Running in parallel

Genetic algorithms are notorious for being [embarrassingly parallel](http://www.wikiwand.com/en/Embarrassingly_parallel). Indeed, most calculations can be run in parallel because they only affect oneindividual. Luckily Go provides good support for parallelism. As some gophers may know, the `math/rand` module can be problematic because there is a global lock the random number generator, the problem is described in this [stackoverflow post](http://stackoverflow.com/questions/14298523/why-does-adding-concurrency-slow-down-this-golang-code). This can be circumvented by providing each GA with it's own generator.
