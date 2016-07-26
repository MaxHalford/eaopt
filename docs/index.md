![logo](img/logo.png)

This is the documentation for [gago](http://gago.readthedocs.io/en/latest/), a genetic algorithm library written in Go. The aim of this project is to make it possible to implement any kind genetic algorithm very quickly. The project is designed with modularity in mind so as to make easy enough to switch genetic operators in a declarative manner. Indeed, every genetic operator can be though of as a brick which can be simply switched with another brick.

The main reason why gago has been implemented with [Go](https://golang.org/) is speed. Genetic algorithms are famous for being [embarrassingly parallel](https://www.wikiwand.com/en/Embarrassingly_parallel) and Go makes it easy to run loops in parallel. What's more, interpreted languages like Python are simply not fast enough for using genetic algorithms in production. Altough not offering as many generic capabilities as other languages, Go is both an extremely quick language and readable one.

A lot of effort has been made to make the code readable and commented. I've also taken the time to provide [some examples](https://github.com/MaxHalford/gago/tree/master/examples) to get started. Feel free to contact me if ever something isn't clear in these examples.
