package completer

import "context"

type Query struct {
	Types  []string
	Term   string
	Locale string
}

type Suggestions struct {
	Items []*Suggestion
	// room for meta
}

type Suggestion struct {
	Slug     string
	Subtitle string
	Title    string
}

type Completer interface {
	Complete(context.Context, *Query) (*Suggestions, error)
}
