package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ivansevryukov1995/github-user-activity/cmd/client"
	"github.com/ivansevryukov1995/github-user-activity/config"
	"github.com/ivansevryukov1995/github-user-activity/internal/model"
	"github.com/ivansevryukov1995/github-user-activity/internal/pkg"
	"github.com/ivansevryukov1995/github-user-activity/internal/pkg/cache"
)

var (
	githubEvents = "https://api.github.com"
	countEvents  = 300
)

func main() {
	clear := flag.Bool("clear-cache", false, "delete file cache")
	flag.Parse()
	args := flag.Args()

	c, err := config.InitConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *clear {
		if err := c.Cache().DeleteCache(); err != nil {
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

	events, hasCache, err := c.Cache().GetCache(user)
	if err != nil {
		if !errors.Is(err, cache.ErrUserNotFound) {
			return fmt.Errorf("failed get cache: %w", err)
		}
		events = make([]model.Event, countEvents)
	}

	if !hasCache {
		clt := &client.GitHubClient{
			URL: "https://api.github.com",
		}

		req := client.GitHubRequest{
			User:        user,
			CountEvents: countEvents,
		}

		events, err = clt.GetEvents(ctx, req)
		if err != nil {
			return fmt.Errorf("failed client: %w", err)
		}

		newUser := &model.User{
			NameUser:       user,
			CreatedAtCache: time.Now().UTC(),
			Events:         events,
		}

		err = c.Cache().SetCache(newUser)
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
