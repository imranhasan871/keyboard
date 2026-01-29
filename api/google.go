package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GoogleInputResponse represents the structure returned by Google Input Tools
// Format is: ["SUCCESS",[["input",["suggestion1","suggestion2",...],...]]]
type GoogleInputResponse [2]interface{}

// FetchSuggestions calls the Google Input Tools API
func FetchSuggestions(text string, lang string) ([]string, error) {
	if lang == "" {
		lang = "bn-t-i0-und" // Default to Bangla
	}

	baseURL := "https://inputtools.google.com/request"
	params := url.Values{}
	params.Add("text", text)
	params.Add("itc", lang)
	params.Add("num", "5") // Number of suggestions
	params.Add("cp", "0")
	params.Add("cs", "1")
	params.Add("ie", "utf-8")
	params.Add("oe", "utf-8")
	params.Add("app", "desktop_client")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GoogleInputResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result[0] != "SUCCESS" {
		return nil, fmt.Errorf("API returned error: %v", result[0])
	}

	// Parsing the dynamic JSON structure safely
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
