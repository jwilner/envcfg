[![Tests](https://github.com/jwilner/envcfg/workflows/tests/badge.svg)](https://github.com/jwilner/envcfg/actions?query=workflow%3Atests+branch%3Amain)
[![Lint](https://github.com/jwilner/envcfg/workflows/lint/badge.svg)](https://github.com/jwilner/envcfg/actions?query=workflow%3Alint+branch%3Amain)
[![GoDoc](https://godoc.org/github.com/jwilner/envcfg?status.svg)](https://godoc.org/github.com/jwilner/envcfg)

# envcfg

A simple and powerful env config library which also provides generated documentation of your configuration values.

```go
package main

import (
    "encoding/json"
    "github.com/jwilner/envcfg"
    "log"
    "os"
)

func main() {
    cfg := envcfg.New()
    time := cfg.Time("EXAMPLE_TIME optional")
    ints := cfg.IntSlice("EXAMPLE_INT_SLICE comma=: base=16 default=a:b:c | Some info about your super important ints")
    
    descriptions, err := cfg.Result()
    if err != nil {
        log.Fatal(err)
    }

    if len(os.Args) > 1 && os.Args[1] == "--help" {
        _ = json.NewEncoder(os.Stderr).Encode(descriptions)
        // Output:
        //[
        //  {"name": "EXAMPLE_TIME", "type": "time.Time", "optional": true},
        //  {
        //      "name": "EXAMPLE_INT_SIZE", 
        //      "type": "[]int64", 
        //      "params": {"comma": ":", "base": 16}, 
        //      "default": [11, 12, 13],
        //      "comment": "Some info about your super important ints"
        //   }
        //]
        os.Exit(1)
    }
    
    _ = time
    _ = ints
}
```
