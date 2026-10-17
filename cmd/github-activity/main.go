package main

import (
	"context"
	"encoding/json"
	"errors"
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

	"github.com/ivansevryukov1995/github-user-activity/config"
	"github.com/ivansevryukov1995/github-user-activity/internal/model"
	"github.com/ivansevryukov1995/github-user-activity/internal/pkg"
	"github.com/ivansevryukov1995/github-user-activity/internal/pkg/cache"
)

var (
	githubEvents = "https://api.github.com/users/%s/events"
	countEvents  = 300
)

func main() {
	clear := flag.Bool("clear-cache", false, "delete file cache")
	flag.Parse()
	args := flag.Args()

	if *clear {
		if err := cache.DeleteCache(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stdout, "Cache clear!")
		os.Exit(0)
	}

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Need one name user-github")
		os.Exit(1)
	}

	c, err := config.InitConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := run(os.Stdout, args[0], c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(out io.Writer, user string, c *config.Config) error {
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

	events, hasCache, err := cache.GetCache(c.Cache(), user)
	if err != nil {
		if !errors.Is(err, cache.ErrUserNotFound) {
			return fmt.Errorf("failed get cache: %w", err)
		}
		events = make([]model.Event, countEvents)
	}

	if !hasCache {
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

		if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
			return fmt.Errorf("failed decode: %w", err)
		}
		newUser := &model.User{
			NameUser:       user,
			CreatedAtCache: time.Now().UTC(),
			Events:         events,
		}
		err = cache.SetCache(c.Cache(), newUser)
		if err != nil {
			return fmt.Errorf("failed set cache: %w", err)
		}
	}

	repoStat := make(map[string]map[model.Action]int, countEvents)

	for _, e := range events {
		event := strings.TrimSuffix(e.Type, "Event")

		if repoStat[e.Repo.Name] == nil {
			repoStat[e.Repo.Name] = make(map[model.Action]int)
		}

		payload := model.ConvertPayloadToAction(event, e)

		repoStat[e.Repo.Name][payload]++
	}

	if hasCache {
		fmt.Fprintln(out, "[CACHE]")
	}

	err = pkg.RenderRepoStatUser(out, repoStat)
	if err != nil {
		return fmt.Errorf("failed render repository statistic user: %w", err)
	}

	return nil
}
