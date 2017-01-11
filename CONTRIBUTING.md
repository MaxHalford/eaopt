## Genetic operators

There are many genetic operators yet to be implemented. Feel free to send pull requests with your implementations. The only requirements are that the implementation respects the existing [naming convention](#naming-convention) and includes a test in the corresponding `*_test.go` file.

## Ideas

- Benchmark against other libraries
- Add more example usage
- Add a way to compare individuals based on their genome to improve speciation
- Make it easier to test models. Also make sure they work as expected.
- Implement operators described in http://www.ppgia.pucpr.br/~alceu/mestrado/aula3/IJBB-41.pdf
- Implement and order mutation operators
- Consider comparing to the best fitness at every individual evaluation, instead of doing a global search after each generation
- Implementing genome type interfaces (for example SliceGenome with At and Set methods)
- Consider using memetic algorithms for local optimization (not well said). For example make it possible to use hill climbing
- Add population and species pointers for individuals
- Implement Partitioning Around Medoids (PAM) for speciating
- Implement Particle Swarm Optimization
- List available operators/models
- http://deap.readthedocs.io/en/master/
- http://pyevolve.sourceforge.net/intro.html#ga-features
- http://www.dmi.unict.it/mpavone/nc-cs/materiale/moscato89.pdf

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

### Variable declaration in returning functions

Avoid declaring returning returning variables inside the function. They are already declared if they are named at the end of the first line of the function.

```go
// Good
func lengthOfList(list []float64) length int {
    length = len(list)
    return
}

// Bad
func lengthOfList(list []float64) int {
    var length int
    length = len(list)
    return length
}
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
| Genetic algorithm | ga           |
| Population        | pop          |


## Parallelism and random number generation caveat

Genetic algorithms are notorious for being [embarrassingly parallel](http://www.wikiwand.com/en/Embarrassingly_parallel). Indeed, most calculations can be run in parallel because they only affect part of the GA. Luckily Go provides good support for parallelism. As some gophers may have encountered, the `math/rand` module can be problematic because there is a global lock attached to the random number generator. The problem is described in this [StackOverflow post](http://stackoverflow.com/questions/14298523/why-does-adding-concurrency-slow-down-this-golang-code). This can be circumvented by providing each population with it's own random number generator.

Talking about parallelism, there is a reason why the populations are run in parallel and not the individuals. First of all for parallelism at an individual level each individual would have to be assigned a new random number generator, which isn't very efficient. Second of all, even though Golang has an efficient concurrency model, spawning routines nonetheless has an overhead. It's simply not worth using a routine for each individual because operations at an individual level are often not time consuming enough.

## Editing the documentation

- The documentation is built with [mkdocs](https://mkdocs.readthedocs.io).
- Each page has an associated markdown file in the `docs/` folder.
- You can `mkdocs serve` to enable live editing of the documentation.
