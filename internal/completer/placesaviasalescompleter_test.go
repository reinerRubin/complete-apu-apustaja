package completer

import (
	"context"
	"testing"
)

// would not work without internet
func TestPlacesavisalesCompleterBasic(t *testing.T) {
	completer, err := NewPlacesaviasales()
	if err != nil {
		t.Fatalf("cant create aviaplaces: %s", err)
	}
	query := &Query{
		Types:  []string{"city", "zoos"},
		Locale: "en",
		Term:   "Moscow",
	}

	ctx := context.Background()
	suggestions, err := completer.Complete(ctx, query)
	if err != nil {
		t.Fatalf("cant get aviaplaces: %s", err)
	}

	if len(suggestions.Items) == 0 {
		t.Fatal("suggestions are empty")
	}
}

func TestPlacesaviasalesQueryBuilderNaive(t *testing.T) {
	pl, err := NewPlacesaviasales()
	if err != nil {
		t.Fatalf("cant create placesavisales completer")
	}

	url, err := pl.buidURL(&Query{
		Types:  []string{"city", "zoos"},
		Locale: "en",
		Term:   "Moscow",
	})
	if err != nil {
		t.Fatalf("cant build URL: %s", err)
	}

	expected := "https://places.aviasales.ru/v2/places.json?%5B%5Dtypes=city&%5B%5Dtypes=zoos&locale=en&term=Moscow"
	if url.String() != expected {
		t.Fatalf("actual: %s, expected: %s", url.String(), expected)
	}
}
