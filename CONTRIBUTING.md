## Ideas

- Add example with penalty
- Add hall of fame field to GA
- Improve tournament selection testing
- Refactor models testings
- Add more context to errors (at least the method/struct name) (https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
- Add more example usage
- Implement operators described in http://www.ppgia.pucpr.br/~alceu/mestrado/aula3/IJBB-41.pdf
- Implement Particle Swarm Optimization
- http://deap.readthedocs.io/en/master/
- http://pyevolve.sourceforge.net/intro.html#ga-features
- http://www.dmi.unict.it/mpavone/nc-cs/materiale/moscato89.pdf
- Serialize with http://labix.org/gobson, maybe

## Code style

### Guidelines

- Keep names short
- Aim for 100 characters per line

### Variable declaration

```go
// Bad
x := 42

// Good
var x = 42
```

## Naming convention

Please inspire yourself from the existing algorithms before implementing, the naming conventions are easy to grasp.


## Parallelism and random number generation caveat

Genetic algorithms are notorious for being [embarrassingly parallel](http://www.wikiwand.com/en/Embarrassingly_parallel). Indeed, most calculations can be run in parallel because they only affect part of the GA. Luckily Go provides good support for parallelism. As some gophers may have encountered, the `math/rand` module can be problematic because there is a global lock attached to the random number generator. The problem is described in this [StackOverflow post](http://stackoverflow.com/questions/14298523/why-does-adding-concurrency-slow-down-this-golang-code). This can be circumvented by providing each population with it's own random number generator.

Talking about parallelism, there is a reason why the populations are run in parallel and not the individuals. First of all for parallelism at an individual level each individual would have to be assigned a new random number generator, which isn't very efficient. Second of all, even though Golang has an efficient concurrency model, spawning routines nonetheless has an overhead. It's simply not worth using a routine for each individual because operations at an individual level are often not time consuming enough.

## Performance

1. `go test -bench . -cpuprofile=cpu.prof`
2. `go tool pprof -pdf gago.test cpu.prof > profile.pdf`
