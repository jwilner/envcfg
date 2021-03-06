package envcfg_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/jwilner/envcfg"
)

func TestConfigure(t *testing.T) {
	tests := []struct {
		name                string
		env                 map[string]string
		configureFunc       func(*envcfg.Cfg) interface{}
		expectedVal         interface{}
		expectedDescription string
		wantErr             string
	}{
		{
			name: "int",
			env:  map[string]string{"a": "1"},
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Int("a")
			},
			expectedVal: int64(1),
			expectedDescription: `{
	"name": "a",
	"type": "int64",
	"optional": false,
	"params": {}
}`,
		},
		{
			name: "int missing",
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Int("a")
			},
			expectedVal: int64(0),
			wantErr:     "a: variable is required",
			expectedDescription: `{
	"name": "a",
	"type": "int64",
	"optional": false,
	"params": {}
}`,
		},
		{
			name: "int default string",
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Int("a default=2")
			},
			expectedVal: int64(2),
			expectedDescription: `{
	"name": "a",
	"type": "int64",
	"optional": false,
	"default": 2,
	"params": {}
}`,
		},
		{
			name: "int default opt",
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Int("a", envcfg.IntDefault(2))
			},
			expectedVal: int64(2),
			expectedDescription: `{
	"name": "a",
	"type": "int64",
	"optional": false,
	"default": 2,
	"params": {}
}`,
		},
		{
			name: "int default opt diff base",
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Int("a base=2 default=010101")
			},
			expectedVal: int64(21),
			expectedDescription: `{
	"name": "a",
	"type": "int64",
	"optional": false,
	"default": 21,
	"params": {
		"base": 2
	}
}`,
		},
		{
			name: "int default opt diff bit_size",
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Int("a base=2 default=010101 bit_size=32")
			},
			expectedVal: int64(21),
			expectedDescription: `{
	"name": "a",
	"type": "int64",
	"optional": false,
	"default": 21,
	"params": {
		"base": 2,
		"bit_size": 32
	}
}`,
		},
		{
			name: "int unspecified base but prefix",
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Int("a default=0xa")
			},
			expectedVal: int64(10),
			expectedDescription: `{
	"name": "a",
	"type": "int64",
	"optional": false,
	"default": 10,
	"params": {}
}`,
		},
		{
			name: "bool",
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Bool("a")
			},
			expectedVal: false,
			wantErr:     "a: variable is required",
			expectedDescription: `{
	"name": "a",
	"type": "bool",
	"optional": false,
	"params": {}
}`,
		},
		{
			name: "bool def true",
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Bool("a default=true")
			},
			expectedVal: true,
			expectedDescription: `{
	"name": "a",
	"type": "bool",
	"optional": false,
	"default": true,
	"params": {}
}`,
		},
		{
			name: "bool def false",
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Bool("a default=false")
			},
			expectedVal: false,
			expectedDescription: `{
	"name": "a",
	"type": "bool",
	"optional": false,
	"default": false,
	"params": {}
}`,
		},
		{
			name: "float",
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Float("a")
			},
			expectedVal: 0.,
			wantErr:     "a: variable is required",
			expectedDescription: `{
	"name": "a",
	"type": "float64",
	"optional": false,
	"params": {}
}`,
		},
		{
			name: "float bit size",
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Float("a bit_size=32")
			},
			env:         map[string]string{"a": ".1234"},
			expectedVal: float64(float32(0.1234)),
			expectedDescription: `{
	"name": "a",
	"type": "float64",
	"optional": false,
	"params": {"bit_size": 32}
}`,
		},
		{
			name: "float default",
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Float("a bit_size=32 default=19")
			},
			expectedVal: float64(float32(19)),
			expectedDescription: `{
	"name": "a",
	"type": "float64",
	"optional": false,
	"default": 19,
	"params": {"bit_size": 32}
}`,
		},
		{
			name: "duration",
			env:  map[string]string{"a": "30m"},
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Duration("a")
			},
			expectedVal: 30 * time.Minute,
			expectedDescription: `{
	"name": "a",
	"type": "time.Duration",
	"optional": false,
	"params": {}
}`,
		},
		{
			name: "duration default",
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Duration("a default=15ms")
			},
			expectedVal: 15 * time.Millisecond,
			expectedDescription: `{
	"name": "a",
	"type": "time.Duration",
	"default": 15000000,
	"optional": false,
	"params": {}
}`,
		},
		{
			name: "time",
			env:  map[string]string{"a": "2019-01-01T02:03:04Z"},
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Time("a")
			},
			expectedVal: time.Date(2019, 01, 01, 02, 03, 04, 0, time.UTC),
			expectedDescription: `{
	"name": "a",
	"type": "time.Time",
	"optional": false,
	"params": {}
}`,
		},
		{
			name: "time custom layout",
			env:  map[string]string{"a": "2019-01-01T02:03Z"},
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Time("a layout=2006-01-02T15:04Z07:00")
			},
			expectedVal: time.Date(2019, 01, 01, 02, 03, 0, 0, time.UTC),
			expectedDescription: `{
	"name": "a",
	"type": "time.Time",
	"optional": false,
	"params": {"layout": "2006-01-02T15:04Z07:00"}
}`,
		},
		{
			name: "int slice",
			env:  map[string]string{"a": "1,2,3,4"},
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.IntSlice("a")
			},
			expectedVal: []int64{1, 2, 3, 4},
			expectedDescription: `{
	"name": "a",
	"type": "[]int64",
	"optional": false,
	"params": {}
}`,
		},
		{
			name: "int slice",
			env:  map[string]string{"a": "1 2 3 a"},
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.IntSlice(`a comma=" " base=16`)
			},
			expectedVal: []int64{1, 2, 3, 10},
			expectedDescription: `{
	"name": "a",
	"type": "[]int64",
	"optional": false,
	"params": {"comma": " ", "base": 16}
}`,
		},
		{
			name: "int slice empty",
			env:  map[string]string{"a": ""},
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.IntSlice(`a comma=" " base=16`)
			},
			expectedVal: []int64{},
			expectedDescription: `{
	"name": "a",
	"type": "[]int64",
	"optional": false,
	"params": {"comma": " ", "base": 16}
}`,
		},
		{
			name: "uint",
			env:  map[string]string{"a": "1"},
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Uint("a")
			},
			expectedVal: uint64(1),
			expectedDescription: `{
	"name": "a",
	"type": "uint64",
	"optional": false,
	"params": {}
}`,
		},
		{
			name: "uint missing",
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Uint("a")
			},
			expectedVal: uint64(0),
			wantErr:     "a: variable is required",
			expectedDescription: `{
	"name": "a",
	"type": "uint64",
	"optional": false,
	"params": {}
}`,
		},
		{
			name: "uint default string",
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Uint("a default=2")
			},
			expectedVal: uint64(2),
			expectedDescription: `{
	"name": "a",
	"type": "uint64",
	"optional": false,
	"default": 2,
	"params": {}
}`,
		},
		{
			name: "uint default opt",
			configureFunc: func(cfg *envcfg.Cfg) interface{} {
				return cfg.Uint("a", envcfg.UintDefault(2))
			},
			expectedVal: uint64(2),
			expectedDescription: `{
	"name": "a",
	"type": "uint64",
	"optional": false,
	"default": 2,
	"params": {}
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newCfg := func() *envcfg.Cfg {
				return envcfg.New(
					envcfg.Panic(false),
					envcfg.EnvFunc(func(k string) (string, bool) {
						v, ok := tt.env[k]
						return v, ok
					}),
				)
			}

			t.Run("configure", func(t *testing.T) {
				c := newCfg()
				val := tt.configureFunc(c)

				var errS string
				if err := c.Err(); err != nil {
					errS = err.Error()
				}
				if errS != tt.wantErr {
					t.Errorf("Configure() error = %q, wantErr %q", errS, tt.wantErr)
				}

				if !reflect.DeepEqual(tt.expectedVal, val) {
					t.Errorf("Configure() want = %v, got %v", tt.expectedVal, val)
				}
			})

			t.Run("describe", func(t *testing.T) {
				c := newCfg()
				_ = tt.configureFunc(c)

				descriptions, _ := c.Result()
				if len(descriptions) != 1 {
					t.Fatalf("Describe() wanted a description, got none")
				}

				var expected, got interface{}
				if err := json.Unmarshal([]byte(tt.expectedDescription), &expected); err != nil {
					t.Fatalf("got a json unmarshal error: %v", err)
				}

				if bs, err := json.Marshal(descriptions[0]); err != nil {
					t.Fatalf("got a json marshal error: %v", err)
				} else if err = json.Unmarshal(bs, &got); err != nil {
					t.Fatalf("got a json unmarshal error: %v", err)
				}

				if !reflect.DeepEqual(expected, got) {
					t.Fatalf("Describe() expected %v, got %v", expected, got)
				}
			})
		})

	}
}
