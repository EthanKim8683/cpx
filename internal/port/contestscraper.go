package port

import (
	"context"

	"github.com/EthanKim8683/cpx/internal/domain"
)

type ContestScraper interface {
	ScrapeContest(ctx context.Context, url string) (*domain.Contest, error)
}
