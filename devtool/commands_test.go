package devtool

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestAddCommand(t *testing.T) {
	AddCommand(defaultCommands()())
	AddCommand(map[string]any{
		"Hello": func(w io.Writer) {
			_, _ = fmt.Fprintf(w, "hello-----\n")
		},
		"custom": map[string]any{
			"sub1": func() {},
			"sub2": func(i int, b bool, s string) error {
				fmt.Println(i, b, s)
				return nil
			},
		},
	})
	help(os.Stdout)

	doCommand(os.Stdout, []string{"Hello"})
	doCommand(os.Stdout, []string{"custom", "sub2", "1", "0", "aaa"})
}

func TestCallFn(t *testing.T) {
	callFn(os.Stdout, func(string, int) error {
		return nil
	}, nil)

	callFn(os.Stdout, func() string {
		return ""
	}, nil)
}
