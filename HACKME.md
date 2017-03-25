## Ideas

- Add more example usage
- Add a way to compare individuals based on their genome to improve speciation
- Make it easier to test models. Also make sure they work as expected.
- Implement operators described in http://www.ppgia.pucpr.br/~alceu/mestrado/aula3/IJBB-41.pdf
- Implement and order mutation operators
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
- Aim for 100 characters per line

### Variable declaration

```go
// Good
x := 42

// Bad
var x = 42
```

The `:=` operator works in the following case.

```go
val1, err := getVal1()
val2, err := getVal2()
```

However the following will not work.

```go
var val1, err = getVal1()
var val2, err = getVal2()
```

That's the only reason why I prefer to stick to `:=` for variable declaration.


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

Please inspire yourself from the existing algorithms before implementing, the naming conventions are easy to grasp.


## Parallelism and random number generation caveat

Genetic algorithms are notorious for being [embarrassingly parallel](http://www.wikiwand.com/en/Embarrassingly_parallel). Indeed, most calculations can be run in parallel because they only affect part of the GA. Luckily Go provides good support for parallelism. As some gophers may have encountered, the `math/rand` module can be problematic because there is a global lock attached to the random number generator. The problem is described in this [StackOverflow post](http://stackoverflow.com/questions/14298523/why-does-adding-concurrency-slow-down-this-golang-code). This can be circumvented by providing each population with it's own random number generator.

Talking about parallelism, there is a reason why the populations are run in parallel and not the individuals. First of all for parallelism at an individual level each individual would have to be assigned a new random number generator, which isn't very efficient. Second of all, even though Golang has an efficient concurrency model, spawning routines nonetheless has an overhead. It's simply not worth using a routine for each individual because operations at an individual level are often not time consuming enough.

## Editing the documentation

- The documentation is built with [mkdocs](https://mkdocs.readthedocs.io).
- Each page has an associated markdown file in the `docs/` folder.
- You can `mkdocs serve` to enable live editing of the documentation.
