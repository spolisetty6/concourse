package commands

import (
	"testing"

	"github.com/concourse/concourse/atc"
)

func TestPipelinesCommandBuildHeader(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		command  PipelinesCommand
		expected []string
	}{
		{
			name:     "default headers",
			command:  PipelinesCommand{},
			expected: []string{"id", "name", "paused", "public", "last updated"},
		},
		{
			name:     "all teams headers",
			command:  PipelinesCommand{All: true},
			expected: []string{"id", "name", "team", "paused", "public", "last updated"},
		},
		{
			name:     "archived headers",
			command:  PipelinesCommand{IncludeArchived: true},
			expected: []string{"id", "name", "paused", "public", "archived", "last updated"},
		},
		{
			name:     "all teams with archived headers",
			command:  PipelinesCommand{All: true, IncludeArchived: true},
			expected: []string{"id", "name", "team", "paused", "public", "archived", "last updated"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.command.buildHeader(); !equalStrings(got, tc.expected) {
				t.Fatalf("unexpected headers\nexpected: %#v\ngot: %#v", tc.expected, got)
			}
		})
	}
}

func TestPipelinesCommandFilterPipelines(t *testing.T) {
	t.Parallel()

	unfiltered := []atc.Pipeline{
		{ID: 1, Name: "active"},
		{ID: 2, Name: "archived", Archived: true},
		{ID: 3, Name: "active-2"},
	}

	t.Run("excludes archived by default", func(t *testing.T) {
		t.Parallel()
		command := PipelinesCommand{}

		filtered := command.filterPipelines(unfiltered)

		if len(filtered) != 2 {
			t.Fatalf("expected 2 pipelines, got %d", len(filtered))
		}
		if filtered[0].ID != 1 || filtered[1].ID != 3 {
			t.Fatalf("unexpected filtered order/contents: %#v", filtered)
		}
	})

	t.Run("includes archived when enabled", func(t *testing.T) {
		t.Parallel()
		command := PipelinesCommand{IncludeArchived: true}

		filtered := command.filterPipelines(unfiltered)

		if len(filtered) != len(unfiltered) {
			t.Fatalf("expected %d pipelines, got %d", len(unfiltered), len(filtered))
		}
		for i := range unfiltered {
			if filtered[i].ID != unfiltered[i].ID {
				t.Fatalf("pipeline at index %d mismatch, expected %d got %d", i, unfiltered[i].ID, filtered[i].ID)
			}
		}
	})

	t.Run("returns empty when all pipelines are archived", func(t *testing.T) {
		t.Parallel()
		command := PipelinesCommand{}

		filtered := command.filterPipelines([]atc.Pipeline{
			{ID: 10, Archived: true},
			{ID: 11, Archived: true},
		})

		if len(filtered) != 0 {
			t.Fatalf("expected 0 pipelines, got %d", len(filtered))
		}
	})

	t.Run("handles empty input", func(t *testing.T) {
		t.Parallel()
		command := PipelinesCommand{}

		filtered := command.filterPipelines(nil)

		if len(filtered) != 0 {
			t.Fatalf("expected 0 pipelines, got %d", len(filtered))
		}
	})
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
