package google

import (
	"encoding/json"
	"fmt"
	"google-input-keyboard/internal/core/domain"
	"net/http"
	"net/url"
)

type GoogleInputGateway struct {
	client *http.Client
}

func NewGoogleInputGateway() domain.SuggestionService {
	return &GoogleInputGateway{
		client: &http.Client{},
	}
}

// Internal structures to parse the dirty JSON from Google
type googleResponse [2]interface{}

func (g *GoogleInputGateway) GetSuggestions(text string, lang string) ([]string, error) {
	if lang == "" {
		lang = "bn-t-i0-und"
	}

	baseURL := "https://inputtools.google.com/request"
	params := url.Values{}
	params.Add("text", text)
	params.Add("itc", lang)
	params.Add("num", "5")
	params.Add("cp", "0")
	params.Add("cs", "1")
	params.Add("ie", "utf-8")
	params.Add("oe", "utf-8")
	// Use 'deskotp' to avoid unrelated metadata
	params.Add("app", "desktop")

	resp, err := g.client.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result googleResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if status, ok := result[0].(string); !ok || status != "SUCCESS" {
		return nil, fmt.Errorf("API returned error")
	}

	data, ok := result[1].([]interface{})
	if !ok || len(data) == 0 {
		return nil, fmt.Errorf("invalid data format")
	}

	firstGroup, ok := data[0].([]interface{})
	if !ok || len(firstGroup) < 2 {
		return nil, fmt.Errorf("invalid group format")
	}

	suggestionList, ok := firstGroup[1].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid suggestions format")
	}

	suggestions := make([]string, 0, len(suggestionList))
	for _, s := range suggestionList {
		if str, ok := s.(string); ok {
			suggestions = append(suggestions, str)
		}
	}

	return suggestions, nil
}
