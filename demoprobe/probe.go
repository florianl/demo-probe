package demoprobe // import "github.com/florianl/demo-probe/demoprobe"

// Config holds the YAML configuration for the demoprobe extension.
//
//	extensions:
//	  demoprobe/my_external_probe:
//	    modulo: 100   # optional; reduces the random payload to [0, modulo)
type Config struct {
	Modulo uint32 `mapstructure:"modulo"`
}
