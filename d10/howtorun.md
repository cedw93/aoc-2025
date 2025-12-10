# Non-standard Run

This solution uses `github.com/draffensperger/golp`, to do this (on a mac)

* `brew install lp_solve`

The run using:

```
CGO_CFLAGS="-I/opt/homebrew/opt/lp_solve/include" CGO_LDFLAGS="-L/opt/homebrew/opt/lp_solve/lib -llpsolve55" go run main.go < input.txt
```