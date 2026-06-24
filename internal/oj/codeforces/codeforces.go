package codeforces

import "github.com/go-rod/rod"

type Codeforces struct {
	browser *rod.Browser
}

func New(b *rod.Browser) *Codeforces {
	return &Codeforces{
		browser: b,
	}
}
