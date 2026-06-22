//go:build integration

package codeforces

import (
	"errors"
	"testing"

	"github.com/EthanKim8683/cpx/internal/domain"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScrapeProblem(t *testing.T) {
	t.Parallel()

	l := launcher.NewUserMode().
		Context(t.Context()).
		Headless(true).
		UserDataDir(t.TempDir())
	wsURL, err := l.Launch()
	require.NoError(t, err)
	t.Cleanup(func() {
		l.Cleanup()
	})

	b := rod.New().
		Context(t.Context()).
		ControlURL(wsURL)
	require.NoError(t, b.Connect())
	t.Cleanup(func() {
		_ = b.Close()
	})

	tests := map[string]struct {
		url     string
		problem *domain.Problem
		err     error
	}{
		"batch": {
			url: "https://codeforces.com/problemset/problem/2158/F1",
			problem: &domain.Problem{
				URL:  "https://codeforces.com/problemset/problem/2158/F1",
				Type: domain.ProblemTypeStdioBatch,
				StdioBatch: &domain.StdioBatch{
					Inputs: []string{
						`3
2
5
7`,
					},
					Outputs: []string{
						`2 2
1 4 4 6 6
4 4 6 6 9 9 4`,
					},
				},
			},
		},
		"interactive": {
			url: "https://codeforces.com/problemset/problem/2196/C2",
			problem: &domain.Problem{
				URL:  "https://codeforces.com/problemset/problem/2196/C2",
				Type: domain.ProblemTypeStdioInteractive,
			},
		},
		"run-twice": {
			url: "https://codeforces.com/contest/2168/problem/A1",
			problem: &domain.Problem{
				URL:  "https://codeforces.com/contest/2168/problem/A1",
				Type: domain.ProblemTypeStdioRunTwice,
			},
		},
		"profile": {
			url: "https://codeforces.com/profile/EthanKim8683",
			err: errors.New("could not determine problem type"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			scraper := New(b)
			problem, err := scraper.ScrapeProblem(t.Context(), test.url)
			if test.err != nil {
				require.EqualError(t, err, test.err.Error())
				assert.Empty(t, problem)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.problem, problem)
			}
		})
	}
}
