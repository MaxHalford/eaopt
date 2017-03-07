# N-queens problem

This example attempts to solve the [n-queens problem](https://developers.google.com/optimization/puzzles/queens). The goal is to find a board layout where n queens do not attack each other. A finite number of solutions exists for each natural integer higher than 3.

Each individual is a board layout represented by a slice of integers. Each element of the slice represents a row of the board and the associated element represents a column. For example the slice `[1, 3, 2, 8, 7, 6, 5, 4]` would be represented by the following board:

    ♕ . . . . . . .
    . . ♕ . . . . .
    . ♕ . . . . . .
    . . . . . . . ♕
    . . . . . . ♕ .
    . . . . . ♕ . .
    . . . . ♕ . . .
    . . . ♕ . . . .

This representation makes it impossible for queens to be on the same row or the same column. For the rows this is automatic because each index in the slice is unique. For the columns we simply have to preserve the uniqueness of each element in the genome; to do so we can use `gago.MutPermuteInt` for mutation and `gago.CrossPMXInt` for crossover.

This example also shows how to print out genomes you define by implementing the `String()` method. In this case by calling `fmt.Println(ga.Best.Genome)` a chess board will be printed out. This can be quite handy for debugging!

**8 queens**

    . . . ♕ . . . .
    . . . . . . ♕ .
    . . . . ♕ . . .
    . . ♕ . . . . .
    ♕ . . . . . . .
    . . . . . ♕ . .
    . . . . . . . ♕
    . ♕ . . . . . .

**12 queens**

    . . . . . . . . . ♕ . .
    ♕ . . . . . . . . . . .
    . . ♕ . . . . . . . . .
    . . . . . . . . . . ♕ .
    . . . . . ♕ . . . . . .
    . . . . . . . ♕ . . . .
    . . . . . . ♕ . . . . .
    . . . . . . . . . . . ♕
    . . . ♕ . . . . . . . .
    . ♕ . . . . . . . . . .
    . . . . ♕ . . . . . . .
    . . . . . . . . ♕ . . .
