# envcfg

A simple and powerful env config library which also provides generated documentation of your configuration values.

```go
package main

import (
    "encoding/json"
    "log"
    "os"
    "time"
    "github.com/jwilner/envcfg"
)

type config struct {
    Time time.Time
    Ints []int64
}

func (c *config) Configure(cfg envcfg.Configurer) {
    c.Time = cfg.Time("EXAMPLE_TIME optional")
    c.Ints = cfg.IntSlice("EXAMPLE_INT_SLICE comma=: base=16 default=a:b:c | Some info about your super important ints")
}

func main() {
    if len(os.Args) > 1 && os.Args[1] == "--help" {
        _ = json.NewEncoder(os.Stderr).Encode(envcfg.Describe(new(config)))
        os.Exit(1)
    }

    var cfg config
    if err := envcfg.Configure(&cfg); err != nil {
        log.Fatal(err)
    }
    _ = cfg
}
```