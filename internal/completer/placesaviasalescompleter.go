package completer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const placesaviasalesBaseURL = "https://places.aviasales.ru/v2/places.json"

type Placesaviasales struct {
	baseURL string
}

func NewPlacesaviasales() (*Placesaviasales, error) {
	return &Placesaviasales{
		baseURL: placesaviasalesBaseURL,
	}, nil
}

func (d *Placesaviasales) Complete(ctx context.Context, q *Query) (*Suggestions, error) {
	url, err := d.buidURL(q)
	if err != nil {
		return nil, fmt.Errorf("cant build url: %s", err)
	}

	log.Printf("GET request to: %s", url.String())
	reqest, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("cant build request: %s", err)
	}

	reqest = reqest.WithContext(ctx)

	client := &http.Client{}
	response, err := client.Do(reqest)
	if err != nil {
		return nil, fmt.Errorf("cant perform request to %s: %s", url, err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("code is not 200: %s", response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("cant read response body: %s", err)
	}

	placesaviasalesSuggestions := make(PlacesaviasalesSuggestions, 0)
	err = json.Unmarshal(body, &placesaviasalesSuggestions)
	if err != nil {
		log.Printf("cant parse: %s", string(body))
		return nil, fmt.Errorf("cant parse body: %s", err)
	}

	suggestions, err := mapPlacesaviasalesSuggestionsToCompleterSuggestions(
		placesaviasalesSuggestions,
	)
	if err != nil {
		return nil, fmt.Errorf("cant map suggestions to completer format: %s", err)
	}
	return suggestions, nil
}

func (d *Placesaviasales) buidURL(q *Query) (*url.URL, error) {
	// TODO optimize me
	url, err := url.Parse(d.baseURL)
	if err != nil {
		return nil, err
	}

	{ // fill query
		urlQuery := url.Query()

		for _, typo := range q.Types {
			urlQuery.Add("[]types", typo)
		}

		if q.Locale != "" {
			urlQuery.Set("locale", q.Locale)
		}
		if q.Term != "" {
			urlQuery.Set("term", q.Term)
		}

		url.RawQuery = urlQuery.Encode()
	}

	return url, nil
}

// TODO: make it less "if else if else whatever"
func mapPlacesaviasalesSuggestionsToCompleterSuggestions(
	asSuggestions PlacesaviasalesSuggestions,
) (*Suggestions, error) {
	suggestions := &Suggestions{
		Items: make([]*Suggestion, 0, len(asSuggestions)),
	}

	for _, asSuggestion := range asSuggestions {
		item := &Suggestion{
			Slug:  asSuggestion.Code,
			Title: asSuggestion.Name,
		}

		// idk, about this; did this "if" according to the small amount examples
		if asSuggestion.Typo == "airport" {
			item.Subtitle = asSuggestion.CityName
		} else {
			item.Subtitle = asSuggestion.CountyName
		}

		suggestions.Items = append(suggestions.Items, item)
	}

	return suggestions, nil
}

type (
	PlacesaviasalesSuggestions []*PlacesaviasalesSuggestion
	PlacesaviasalesSuggestion  struct {
		Typo       string `json:"type"`
		Code       string `json:"code"`
		CountyName string `json:"country_name"`
		Name       string `json:"name"`
		CityName   string `json:"city_name"`
	}
)
