package envcfg_test

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jwilner/envcfg"
)

func ExampleCfg_Int() {
	os.Setenv("EXAMPLE_HEX_INT", "abcdef")
	c := envcfg.New()
	val := c.Int("EXAMPLE_HEX_INT base=16 default=1f | A hex int configuration value")
	desc, err := c.Result()
	fmt.Printf("%v\n", val)
	fmt.Printf("%v\n", err)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	_ = enc.Encode(desc[0])
	// Output:
	// 11259375
	// <nil>
	//{
	//     "name": "EXAMPLE_HEX_INT",
	//     "type": "int64",
	//     "optional": false,
	//     "default": 31,
	//     "params": {
	//         "base": 16
	//     },
	//     "comment": "A hex int configuration value"
	//}
}
