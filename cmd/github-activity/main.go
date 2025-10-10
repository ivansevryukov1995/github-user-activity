package main

import (
	"context"
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

type Event struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Actor struct {
		Login string `json:"login"`
	} `json:"actor"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
}

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

	client := &http.Client{
		Timeout: time.Duration(time.Second * 10),
	}

	urlEvents := fmt.Sprintf(githubEvents, args[0])

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, urlEvents, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Status code: %v", resp.StatusCode))
	}

	var ev []Event
	if err := json.NewDecoder(resp.Body).Decode(&ev); err != nil {
		return err
	}

	for _, e := range ev {
		fmt.Fprintln(out, e)
	}

	return nil
}
