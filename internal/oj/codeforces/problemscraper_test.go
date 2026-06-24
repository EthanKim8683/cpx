package codeforces

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/EthanKim8683/cpx/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseProblem(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		path    string
		problem *domain.Problem
		err     error
	}{
		"stdiobatch": {
			path: filepath.Join("testdata", "problemscraper", "stdiobatch.html"),
			problem: &domain.Problem{
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
		"stdiointeractive": {
			path: filepath.Join("testdata", "problemscraper", "stdiointeractive.html"),
			problem: &domain.Problem{
				Type: domain.ProblemTypeStdioInteractive,
			},
		},
		"stdioruntwice": {
			path: filepath.Join("testdata", "problemscraper", "stdioruntwice.html"),
			problem: &domain.Problem{
				Type: domain.ProblemTypeStdioRunTwice,
			},
		},
		"profile": {
			path: filepath.Join("testdata", "problemscraper", "profile.html"),
			err:  errors.New("could not determine problem type"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			d := readDocument(t, test.path)
			problem, err := parseProblem(d)
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
