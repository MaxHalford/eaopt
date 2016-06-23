![logo](img/logo.png)

Apple
:   Pomaceous fruit of plants of the genus Malus in
    the family Rosaceae.

!!! tip "Tip"
    ...

!!! caution "Caution"
    ...

!!! note "Note"
    ...

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

For full documentation visit [mkdocs.org](http://mkdocs.org).

When $a \ne 0$, there are two solutions to \(ax^2 + bx + c = 0\) and they are
$$x = {-b \pm \sqrt{b^2-4ac} \over 2a}.$$

## Commands

- `mkdocs new [dir-name]` - Create a new project.
- `mkdocs serve` - Start the live-reloading docs server.
- `mkdocs build` - Build the documentation site.
- `mkdocs help` - Print this help message.

## Project layout

    mkdocs.yml    # The configuration file.
    docs/
        index.md  # The documentation homepage.
        ...       # Other markdown pages, images and other files.
