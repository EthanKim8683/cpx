package domain

type Contest struct {
	URL      string
	Problems map[string]Problem
}
