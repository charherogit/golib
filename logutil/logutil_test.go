package logutil

import (
	"github.com/caarlos0/env/v8"
	"github.com/sirupsen/logrus"
	"os"

	"testing"
)

func TestConfig(t *testing.T) {
	if err := os.Setenv("LogLevel", logrus.ErrorLevel.String()); err != nil {
		t.Fatal(err)
	}
	c := &LogConfig{}
	if err := env.Parse(c); err != nil {
		t.Fatal(err)
	}
	if c.Level != logrus.ErrorLevel {
		t.Fatalf("level error: %v", c.Level)
	}
}

func TestAAA(t *testing.T) {
	Debug("hello world")
	yyy()
}

func yyy() {
	Outer().Debug("dlrow olleh")
}
