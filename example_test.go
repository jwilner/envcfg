package envcfg_test

import (
	"fmt"
	"github.com/jwilner/envcfg"
	"os"
	"time"
)

func ExampleCfg() {
	type coolMan struct {
		Hi  int64
		Bye string
		OK  time.Time
	}

	e := envcfg.New()
	c := coolMan{
		Hi:  e.Int("HI optional"),
		Bye: e.String("BYE optional"),
		OK:  e.Time("OK layout=\"" + time.RFC822 + "\""),
	}
	fmt.Printf("%#v\n", c)
	fmt.Println(e.Err())
	// Output:
	// envcfg_test.coolMan{Hi:0, Bye:"", OK:time.Time{wall:0x0, ext:0, loc:(*time.Location)(nil)}}
	// OK: variable is required
}

func ExampleCfg_Has() {
	_ = os.Setenv("EXAMPLE_CFG_HAS", "1234")
	ok := envcfg.New().Has("EXAMPLE_CFG_HAS")
	fmt.Printf("%v", ok)
	// Output:
	// true
}

func ExampleCfg_Has_Fail() {
	ok := envcfg.New().Has("EXAMPLE_CFG_HAS_FAIL")
	fmt.Printf("%v", ok)
	// Output:
	// false
}



func ExampleCfg_Int() {
	_ = os.Setenv("EXAMPLE_CFG_INT", "129")
	val := envcfg.New().Int("EXAMPLE_CFG_INT")
	fmt.Printf("%v", val)
	// Output:
	// 129
}

func ExampleCfg_Int_Fail() {
	c := envcfg.New()
	val := c.Int("EXAMPLE_CFG_INT_FAIL")
	fmt.Printf("%v\n", val)
	fmt.Printf("%v\n", c.Err())
	// Output:
	// 0
	// EXAMPLE_CFG_INT_FAIL: variable is required
}

func ExampleCfg_Int_Optional() {
	c := envcfg.New()
	val := c.Int("EXAMPLE_CFG_INT_OPTIONAL optional")
	fmt.Printf("%v\n", val)
	fmt.Printf("%v\n", c.Err())
	// Output:
	// 0
	// <nil>
}

func ExampleCfg_Int_Default() {
	val := envcfg.New().Int("EXAMPLE_CFG_INT_DEFAULT default=23")
	fmt.Printf("%v", val)
	// Output:
	// 23
}

func ExampleCfg_StringSlice() {
	_ = os.Setenv("EXAMPLE_CFG_STRINGSLICE", "a|b|c")
	val := envcfg.New().StringSlice(`EXAMPLE_CFG_STRINGSLICE comma=|`)
	fmt.Printf("%v", val)
	// Output:
	// [a b c]
}

func ExampleCfg_IntSlice() {
	_ = os.Setenv("EXAMPLE_CFG_INTSLICE", "a|b|c")
	val := envcfg.New().IntSlice(`EXAMPLE_CFG_INTSLICE comma=| base=16`)
	fmt.Printf("%v", val)
	// Output:
	// [10 11 12]
}
