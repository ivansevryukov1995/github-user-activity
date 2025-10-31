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
	"os/signal"
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

	ctx := context.Background()

	if err := run(ctx, args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, out io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if len(args) != 1 {
		return errors.New("Ukazhite odno imya user-github")
	}

	client := &http.Client{
		Timeout: time.Duration(time.Second * 10),
	}

	urlEvents := fmt.Sprintf(githubEvents, args[0])

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlEvents, nil)
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
		fmt.Fprintln(out, e.Type)
	}

	return nil
}
