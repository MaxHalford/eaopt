## Genetic operators

There are many genetic operators yet to be implemented. Feel free to send pull requests with your implementations. The only requirements are that the genetic operator respects the existing [naming convention](#naming-convention) and includes a test in the corresponding `*_test.go` file.

## Roadmap

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

### Guidelines

- Keep names short
- Aim for 80 characters per line

### Variable declaration

```go
// Good
var x = 42

// Bad
x := 42
```

### Consecutive variable declarations

```go
// Good
var (
    a = 1
    b = 2
    c = 3
)

// Bad
var a = 1
var b = 2
var c = 3
```

## Naming convention

### Genetic operators

Each genetic operator has a prefix and a suffix to easily identify it.

#### Prefix

The name of the operator begins with an abreviation of the kind of operator being implemented:

| Operator    | Abbreviation |
|-------------|--------------|
| Clusterer   | Clu          |
| Crossover   | Cross        |
| Initializer | Init         |
| Model       | Mod          |
| Migrator    | Mig          |
| Mutator     | Mut          |
| Selector    | Sel          |

#### Suffix

Finally the name of the operator ends with a letter indicating on what kind of genomes the operator works:

| Type    | Abbreviation |
|---------|--------------|
| float64 | F            |
| string  | S            |

No suffix indicates that the genetic operator works on any kind of genome, regardless of the underlying type.


### Shortnames

Along with the genetic operators prefixes, other shortnames are used in the code:

| Name              | Abbreviation |
|-------------------|--------------|
| Individual        | indi         |
| Individuals       | indis        |
| Fitness function  | ff           |
| Genetic algorithm | GA           |
| Population        | pop          |


## Parallelism and random number generation caveat

Genetic algorithms are notorious for being [embarrassingly parallel](http://www.wikiwand.com/en/Embarrassingly_parallel). Indeed, most calculations can be run in parallel because they only affect part of the GA. Luckily Go provides good support for parallelism. As some gophers may have encountered, the `math/rand` module can be problematic because there is a global lock attached to the random number generator. The problem is described in this [StackOverflow post](http://stackoverflow.com/questions/14298523/why-does-adding-concurrency-slow-down-this-golang-code). This can be circumvented by providing each population with it's own random number generator.

## Editing the documentation

- The documentation is built with [mkdocs](https://mkdocs.readthedocs.io).
- You can `mkdocs serve` to enable live editing of the documentation.
