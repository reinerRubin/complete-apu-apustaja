package completer

import (
	"context"
)

type DummyCompleter struct {
	Body *Suggestions
}

func NewDummyCompleter() *DummyCompleter {
	return &DummyCompleter{
		Body: &Suggestions{
			Items: []*Suggestion{
				&Suggestion{
					Slug:     "MOW",
					Subtitle: "Russia",
					Title:    "Moscow",
				},
			},
		},
	}
}

func (d *DummyCompleter) Complete(ctx context.Context, q *Query) (*Suggestions, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return d.Body, nil
}
