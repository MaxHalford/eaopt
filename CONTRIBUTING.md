## Ideas

- Implement [differential evolution](https://www.wikiwand.com/en/Differential_evolution)
- Implement [CMAES](https://www.wikiwand.com/en/CMA-ES)
- Implement Particle Swarm Optimization
- Improve tournament selection testing
- Refactor models testings
- Add more context to errors (at least the method/struct name) (https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
- Implement operators described in http://www.ppgia.pucpr.br/~alceu/mestrado/aula3/IJBB-41.pdf
- http://deap.readthedocs.io/en/master/
- http://pyevolve.sourceforge.net/intro.html#ga-features
- http://www.dmi.unict.it/mpavone/nc-cs/materiale/moscato89.pdf

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

## Performance

1. `go test -bench . -cpuprofile=cpu.prof`
2. `go tool pprof -pdf gago.test cpu.prof > profile.pdf`
