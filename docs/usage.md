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


## Parameters

To modify the behavior off the GA, you can change the `gago.GA` struct before running `ga.Initialize`. You can either instantiate a new `gago.GA` or use a predefined one from the `configuration.go` file.

| Variable in the code   | Type                      | Description                                                      |
|------------------------|---------------------------|------------------------------------------------------------------|
| `NbrPopulations`              | `int`                     | Number of Populations in the GA                               |
| `NbrIndividuals`        | `int`                     | Number of individuals in each population                              |
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

## Advice

- Wrap multiple mutators into a single `Mutator` `struct` if you wish to apply multiple mutators, for an example see the [TSP preset](presets/tsp.go).
- Don't hesitate to add more populations if you have a multi-core machine, the overhead is very small.
- Consider the fact that most of the computation is for evaluating the fitness function.
- Increasing the number of selected parents (`NbParents`) during selection usually increases the converrngce rate (which is not necessarily good, but is sometimes desired).
- Increasing the number of individuals per population (`NbrIndividuals`) adds variety to the genetic algorithm, however it is more costly.
- You can access the GA's `duration` attribute or implement your own stopwatch to enhance the GA for a fixed duration.
