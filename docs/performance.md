```sh
go test -bench . -cpuprofile=cpu.prof
go tool pprof gago.test cpu.prof
```
