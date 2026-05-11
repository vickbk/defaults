package defaults

import "errors"

// defaults.Apply executes a series of modifier functions (default.Applier) on a target struct pointer.
// It is designed to safely handle default values using struct initialization and configuration overrides.
// The function takes a pointer to a struct of type T and a variadic list of Applier functions that modify the struct.
// Each Applier function is responsible for applying a specific configuration or validation to the struct.
// If any Applier returns an error, Apply collects these errors and returns them as a single error using errors.Join.
// The function ensures that the target struct is not nil before applying any modifications, returning an error if it is.
// This design allows for flexible and composable configuration of structs while providing robust error handling.
//
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

	var errs []error

	for _, initializer := range initializers {
		if initializer != nil {
			if err := initializer(target); err != nil {
				if errs == nil {
					errs = make([]error, 0, len(initializers))
				}
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return target, errors.Join(errs...)
	}

	return target, nil
}
