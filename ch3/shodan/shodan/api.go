package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIInfo ...
// Data structure used to unroll JSON data returned by /api-info API endpoint.
// Holds basic Shodan user data, eg. Query Credits and Scan Credits.
type APIInfo struct {
	QueryCredits int    `json:"query_credits"`
	ScanCredits  int    `json:"scan_credits"`
	Telnet       bool   `json:"telnet"`
	Plan         string `json:"plan"`
	HTTPS        bool   `json:"https"`
	Unlocked     bool   `json:"unlocked"`
}

// APIInfo ...
// Queries the /api-info API for user information.
// Retrieves basic Shodan user data, eg. Query Credits and Scan Credits.
func (s *Client) APIInfo() (*APIInfo, error) {
	res, err := http.Get(fmt.Sprintf("%s/api-info?key=%s", BaseURL, s.apiKey))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var ret APIInfo
	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
