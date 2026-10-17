package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	binName = "github-activity"
)

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running test...")
	result := m.Run()

	fmt.Println("Cleaning test...")
	os.Remove(binName)

	os.Exit(result)
}

func MockGitHubServerExistUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[{"id":"1", "type":"PushEvent", "actor":{"login":"ivansevryukov1995"}}]`))
}

func MockGitHubServerNotExistUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "Not Found"}`))
}

func TestGithubActivityCLI(t *testing.T) {
	testCases := []struct {
		name   string
		server func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:   "CheckEventsExistUser",
			server: MockGitHubServerExistUser,
		},
		{
			name:   "CheckEventsNotExistUser",
			server: MockGitHubServerNotExistUser,
		},
	}
	mockGitHub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/users/ivansevryukov1995/events") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id":"1", "type":"PushEvent", "actor":{"login":"ivansevryukov1995"}}]`))
			return
		}

		if strings.Contains(r.URL.Path, "/users/Xivansevryukov1995X/events") {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "Not Found"}`))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}))
	defer mockGitHub.Close()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdPath := filepath.Join(dir, binName)
	t.Run("CheckEventsExistUser", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "ivansevryukov1995")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("CheckEventsNotExistUser", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "Xivansevryukov1995X")
		if err := cmd.Run(); err == nil {
			t.Fatal(err)
		}
	})
	t.Run("ClearCache", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-clear-cache")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})
}
