package pkg

import (
	"fmt"
	"io"

	"github.com/ivansevryukov1995/github-user-activity/internal/model"
)

func RenderRepoStatUser(out io.Writer, repoStat map[string]map[model.Action]int) error {
	if len(repoStat) == 0 {
		_, err := fmt.Fprintln(out, "No events to display.")
		if err != nil {
			return fmt.Errorf("failed write in out: %w", err)
		}
		return nil
	}

	_, err := fmt.Fprintln(out, "Output:")
	if err != nil {
		return fmt.Errorf("failed write in out: %w", err)
	}

	if len(repoStat) != 0 {
		for repo, actions := range repoStat {
			for object, count := range actions {
				if count == 1 {
					_, err := fmt.Fprintln(out, "-", object.Action, "a new", object.RefType, "in", repo)
					if err != nil {
						return fmt.Errorf("failed write in out: %w", err)
					}
					continue
				}
				_, err := fmt.Fprintln(out, "-", object.Action, count, object.RefType+"s", "in", repo)
				if err != nil {
					return fmt.Errorf("failed write in out: %w", err)
				}
			}
		}
	}

	return nil
}
