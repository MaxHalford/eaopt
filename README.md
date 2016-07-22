<div align="center">
  <!-- Logo -->
   <img src="docs/img/logo.png" alt="logo"/>
</div>

<div align="center">
  <!-- License -->
  <a href="https://opensource.org/licenses/MIT">
    <img src="http://img.shields.io/:license-mit-ff69b4.svg?style=flat-square" alt="mit"/>
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

<div align="center">For an introduction, example usage, contributing guidelines, please refer to the <a href="http://gago.readthedocs.io/"><b>documentation</b></a>.</div>

## Quick start

It's relatively easy to start using gago by using a [preset](http://gago.readthedocs.io/en/latest/presets/).

```go
package main

import (
    "fmt"
    m "math"

    "github.com/MaxHalford/gago/presets"
)

// Sphere function minimum is 0 reached in (0, ..., 0).
func sphere(X []float64) float64 {
    sum := 0.0
    for _, x := range X {
        sum += m.Pow(x, 2)
    }
    return sum
}

func main() {
    // Instantiate a GA with 2 variables and the fitness function
    var ga = presets.Float(2, sphere)
    ga.Initialize()
    // Enhancement
    for i := 0; i < 1000; i++ {
      ga.Enhance()
    }
    // Display the best obtained solution
    fmt.Printf("The best obtained solution is %f\n", ga.Best.Fitness)
}
```

A preset is simply a genetic algorithm configuration. It's unlikely that a preset will find an optimal solution as is. Presets should be considered as starting points and should be tuned for specific problems.

```go
// Float returns a configuration for minimizing continuous mathematical
// functions with a given number of variables.
func Float(n int, function func([]float64) float64) gago.GA {
    return gago.GA{
        NbrPopulations: 2,
        NbrIndividuals: 30,
        NbrGenes:       n,
        Ff: gago.Float64Function{
            Image: function,
        },
        Initializer: gago.InitUniformF{
            Lower: -1,
            Upper: 1,
        },
        Model: gago.ModGenerational{
            Selector: gago.SelTournament{
                NbParticipants: 3,
            },
            Crossover: gago.CrossUniformF{},
            Mutator: gago.MutNormalF{
                Rate: 0.5,
                Std:  3,
            },
            MutRate: 0.5,
        },
        Migrator:     gago.MigShuffle{},
        MigFrequency: 10,
    }
}
```

## Alternatives

- [GeneticGo](https://github.com/handcraftsman/GeneticGo)
- [goga](https://github.com/tomcraven/goga)
- [go-galib](https://github.com/thoj/go-galib)
