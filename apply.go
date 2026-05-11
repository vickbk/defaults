package defaults

import "errors"

// Apply executes a series of modifier functions (Applier) on a target struct pointer.
// It is designed to safely handle struct initialization and configuration overrides.
//
// Features:
// - Returns the target pointer to allow for inline initialization.
//
// If the target pointer is nil, Apply returns a nil pointer and an error.
// Example usage:
//
//		type Config struct {
//		    Port int
//		    Host string
//		}
//
//		func WithPort(port int) Applier[Config] {
//		    return func(c *Config) error {
//		        if port <= 0 {
//		            return fmt.Errorf("invalid port: %d", port)
//		        }
//		        c.Port = port
//		        return nil
//		    }
//		}
//
//	 func CheckRequredFields(c *Config) error {
//	     if c.Host == "" {
//	         return fmt.Errorf("Host field is empty")
//	     }
//	     if c.Port == 0 {
//	         return fmt.Errorf("Port field is zero")
//	     }
//	     return nil
//	 }
//
//	 func MyFunction(setters ...Applier[Config]) {
//	     cfg, err := Apply(&Config{Host: "localhost", Port: 8080}, setters...)
//	     if err != nil {
//	         // handle error
//	     }
//	     // use cfg
//	 }
//
//	 func main() {
//	     MyFunction(WithPort(9000), CheckRequredFields)
//	 }
func Apply[T any](target *T, initializers ...Applier[T]) (*T, error) {
	if target == nil {
		return nil, errors.New("target cannot be nil")
	}

	var errs []error = make([]error, 0, len(initializers))

	for _, initializer := range initializers {
		if initializer != nil {
			if err := initializer(target); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return target, errors.Join(errs...)
	}

	return target, nil
}
