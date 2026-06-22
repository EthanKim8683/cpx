package port

import (
	"context"

	"github.com/EthanKim8683/cpx/internal/domain"
)

type ProblemScraper interface {
	ScrapeProblem(ctx context.Context, url string) (*domain.Problem, error)
}
