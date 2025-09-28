package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var githubEvents = "https://api.github.com/users/%s/events"

func main() {
	flag.Parse()
	args := flag.Args()

	if err := run(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run(args []string, out io.Writer) error {
	if len(args) != 1 {
		return errors.New("Ukazhite odno imya user-github")
	}

	client := http.Client{
		Timeout: time.Duration(time.Second * 10),
	}

	urlEvents := fmt.Sprintf(githubEvents, args[0])

	res, err := client.Get(urlEvents)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Status code: %v", res.StatusCode))
	}

	data := []byte{}
	if err := json.Unmarshal(data, res.Body); err != nil {
		return err
	}

	fmt.Fprintln(out, data)

	return nil
}
