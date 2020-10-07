/*
Package envcfg is a simple and opinionated library for parsing environment variables and documenting their usage.

The central type is envcfg.Cfg, which can be constructed with envcfg.New. Cfg has a series of parse methods. Each parse
method expects the name of the environment variable and a variety of optional parameters.

The optional parameters can be specified either in the primary string following the comment or with type safe options.
For example, the following two invocations are functionally identical:

	cfg.Int("MY_EXAMPLE_INT", envcfg.IntDefault(10), envcfg.IntBase(16), envcfg.Comment("A really cool int"))

	cfg.Int("MY_EXAMPLE_INT default=a base=16 | A really cool int")

After parsing, errors can be checked with `Err`:

	if err := cfg.Err(); err != nil {
		log.Fatalf("failed loading config: %v", err)
	}

Describe of the config interface can also be printed (e.g. as JSON):

	json.NewEncoder(os.Stdout).Encode(cfg.Describe())

*/
package envcfg
