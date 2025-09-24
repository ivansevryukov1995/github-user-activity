package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if err := run(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run(username []string, out io.Writer) error {
	if len(username) == 0 {
		return errors.New("Ukazhite imya user-github")
	}

	fmt.Fprintln(out, username)

	return nil
}
