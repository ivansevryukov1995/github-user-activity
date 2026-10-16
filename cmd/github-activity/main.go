package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ivansevryukov1995/github-user-activity/internal/model"
	"github.com/ivansevryukov1995/github-user-activity/internal/pkg"
)

var (
	githubEvents = "https://api.github.com/users/%s/events"
	countEvents  = 300
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("Need one name user-github")
	}

	if err := run(args[0], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(user string, out io.Writer) error {
	// graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v, shutting down gracefully...", sig)
		cancel()
	}()

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	urlEvents := fmt.Sprintf(githubEvents, user)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlEvents, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %v", resp.StatusCode)
	}

	events := make([]model.Event, countEvents)
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return fmt.Errorf("failed decode: %w", err)
	}

	repoStat := make(map[string]map[model.Payload]int, countEvents)

	for _, e := range events {
		event := strings.TrimSuffix(e.Type, "Event")

		if repoStat[e.Repo.Name] == nil {
			repoStat[e.Repo.Name] = make(map[model.Payload]int)
		}

		payload := model.ConvertPayload(event, e)

		repoStat[e.Repo.Name][payload]++
	}

	err = pkg.RenderRepoStatUser(out, repoStat)
	if err != nil {
		return fmt.Errorf("failed render repository statistic user: %w", err)
	}

	return nil
}
