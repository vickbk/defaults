package defaults

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

type applyConfig struct {
	Port int
	Host string
}

func withPort(port int) Applier[applyConfig] {
	return func(c *applyConfig) error {
		if port <= 0 {
			return fmt.Errorf("invalid port: %d", port)
		}
		c.Port = port
		return nil
	}
}

func withHost(host string) Applier[applyConfig] {
	return func(c *applyConfig) error {
		if host == "" {
			return errors.New("missing host")
		}
		c.Host = host
		return nil
	}
}

func failingInitializer(message string) Applier[applyConfig] {
	return func(*applyConfig) error {
		return fmt.Errorf(message)
	}
}

func TestApply(t *testing.T) {
	tests := []struct {
		name              string
		target            *applyConfig
		initializers     []Applier[applyConfig]
		wantPort          int
		wantHost          string
		wantErr           bool
		wantErrContains   []string
	}{
		{
			name: "Success applies multiple initializers",
			target: &applyConfig{Port: 80, Host: "localhost"},
			initializers: []Applier[applyConfig]{
				withPort(9000),
				withHost("example.com"),
			},
			wantPort: 9000,
			wantHost: "example.com",
			wantErr:  false,
		},
		{
			name: "Nil target returns error",
			target: nil,
			initializers: []Applier[applyConfig]{
				withPort(9000),
			},
			wantErr:         true,
			wantErrContains: []string{"target cannot be nil"},
		},
		{
			name: "Nil initializer is ignored",
			target: &applyConfig{Port: 80, Host: "localhost"},
			initializers: []Applier[applyConfig]{
				withPort(9000),
				nil,
				withHost("example.com"),
			},
			wantPort: 9000,
			wantHost: "example.com",
			wantErr:  false,
		},
		{
			name: "Multiple errors are joined",
			target: &applyConfig{Port: 80, Host: "localhost"},
			initializers: []Applier[applyConfig]{
				failingInitializer("first failure"),
				failingInitializer("second failure"),
			},
			wantErr:         true,
			wantErrContains: []string{"first failure", "second failure"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Apply(tt.target, tt.initializers...)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				for _, wantText := range tt.wantErrContains {
					if !strings.Contains(err.Error(), wantText) {
						t.Fatalf("expected error to contain %q, got %q", wantText, err.Error())
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result == nil {
				t.Fatal("expected non-nil result")
			}
			if result.Port != tt.wantPort {
				t.Fatalf("expected port %d, got %d", tt.wantPort, result.Port)
			}
			if result.Host != tt.wantHost {
				t.Fatalf("expected host %q, got %q", tt.wantHost, result.Host)
			}
		})
	}
}

func BenchmarkApplySuccess(b *testing.B) {
	target := &applyConfig{Port: 80, Host: "localhost"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Apply(target, withPort(9000))
	}
}

func BenchmarkApplyWithErrors(b *testing.B) {
	target := &applyConfig{Port: 80, Host: "localhost"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Apply(target, failingInitializer("first"), failingInitializer("second"))
	}
}
