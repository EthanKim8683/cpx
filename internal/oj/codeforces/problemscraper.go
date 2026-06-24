package codeforces

import (
	"context"
	"fmt"
	"strings"

	"github.com/EthanKim8683/cpx/internal/domain"
	"github.com/EthanKim8683/cpx/internal/port"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod/lib/proto"
)

func parseProblemType(d *goquery.Document) (domain.ProblemType, error) {
	var (
		hasFirstRun    = false
		hasSecondRun   = false
		hasInteraction = false
		hasInput       = false
		hasOutput      = false
	)
	d.Find("div.section-title, span.tex-font-style-bf").
		Each(func(_ int, s *goquery.Selection) {
			switch strings.TrimSpace(s.Text()) {
			case "First Run":
				hasFirstRun = true
			case "Second Run":
				hasSecondRun = true
			case "Interaction":
				hasInteraction = true
			case "Input":
				hasInput = true
			case "Output":
				hasOutput = true
			}
		})

	switch {
	case hasFirstRun && hasSecondRun:
		return domain.ProblemTypeStdioRunTwice, nil
	case hasInteraction:
		return domain.ProblemTypeStdioInteractive, nil
	case hasInput && hasOutput:
		return domain.ProblemTypeStdioBatch, nil
	default:
		return domain.ProblemTypeUnspecified, fmt.Errorf("could not determine problem type")
	}
}

func innerText(s *goquery.Selection) string {
	divs := s.Find("div")
	if divs.Length() == 0 {
		return strings.TrimSpace(s.Text())
	}

	var sb strings.Builder
	divs.Each(func(_ int, d *goquery.Selection) {
		sb.WriteString(d.Text())
		sb.WriteString("\n")
	})
	return strings.TrimSpace(sb.String())
}

func parseStdioBatch(d *goquery.Document) *domain.StdioBatch {
	inputs := d.Find("div.input pre").
		Map(func(_ int, s *goquery.Selection) string {
			return innerText(s)
		})
	outputs := d.Find("div.output pre").
		Map(func(_ int, s *goquery.Selection) string {
			return innerText(s)
		})
	return &domain.StdioBatch{
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func parseProblem(d *goquery.Document) (*domain.Problem, error) {
	problemType, err := parseProblemType(d)
	if err != nil {
		return nil, err
	}

	problem := &domain.Problem{
		Type: problemType,
	}
	switch problemType {
	case domain.ProblemTypeStdioBatch:
		problem.StdioBatch = parseStdioBatch(d)
	case domain.ProblemTypeStdioInteractive:
	case domain.ProblemTypeStdioRunTwice:
	}
	return problem, nil
}

func (c *Codeforces) ScrapeProblem(ctx context.Context, url string) (*domain.Problem, error) {
	b := c.browser.Context(ctx)

	page, err := b.Page(proto.TargetCreateTarget{
		URL: url,
	})
	if err != nil {
		return nil, err
	}
	defer page.Close()

	if err := page.WaitLoad(); err != nil {
		return nil, err
	}

	html, err := page.HTML()
	if err != nil {
		return nil, err
	}

	d, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	problem, err := parseProblem(d)
	if err != nil {
		return nil, err
	}
	problem.URL = url
	return problem, nil
}

var _ port.ProblemScraper = (*Codeforces)(nil)
