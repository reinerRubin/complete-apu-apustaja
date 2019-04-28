package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/reinerRubin/complete-apu-apustaja/internal/completer"
)

func NewCompleterHandler(completer completer.Completer) *CompleterHandler {
	return &CompleterHandler{
		completer: completer,
	}
}

type CompleterHandler struct {
	completer completer.Completer
}

type completerResponse struct {
	suggestions *completer.Suggestions
	err         error
}

func (h *CompleterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	query, err := parseCompleteQueryParams(r.URL)
	if err != nil {
		WriteError(w, NewErrorResponse(err.Error()))
		return
	}

	result, err := h.complete(r.Context(), query)
	if err != nil {
		WriteError(w, NewErrorResponse(err.Error()))
		return
	}

	// TODO: add a mapper from completer models to transport
	WriteResult(w, result.Items)
}

func (h *CompleterHandler) complete(
	ctx context.Context,
	query *completer.Query,
) (*completer.Suggestions, error) {
	completeResponseChan := make(chan *completerResponse)
	go func() {
		suggestions, err := h.completer.Complete(ctx, query)
		completeResponseChan <- &completerResponse{
			err:         err,
			suggestions: suggestions,
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-completeResponseChan:
		if result.err != nil {
			return nil, result.err
		}

		return result.suggestions, nil
	}
}

func parseCompleteQueryParams(url *url.URL) (*completer.Query, error) {
	query := url.Query()
	types := query["types[]"]
	if len(types) == 0 {
		return nil, fmt.Errorf("types is not provided")
	}

	term := query.Get("term")
	if term == "" {
		return nil, fmt.Errorf("term is not provided")
	}

	locale := query.Get("locale")
	if locale == "" {
		locale = "en"
	}

	return &completer.Query{
		Types:  types,
		Term:   term,
		Locale: locale,
	}, nil
}
