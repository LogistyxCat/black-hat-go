package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Profile ...
// Used to hold queried data from /account/profile API endpoint.
// Populated by the GetProfileInfo() method.
type Profile struct {
	Member        bool   `json:"member"`
	ExportCredits int    `json:"credits"`
	DisplayName   string `json:"display_name"`
	CreatedDate   string `json:"created"`
}

// GetProfileInfo ...
// Queries the /account/profile API endpoint for Shodan user profile info.
// Unwraps response into Profile object and returns pointer to object and an appropriate error.
func (s *Client) GetProfileInfo() (*Profile, error) {
	res, err := http.Get(fmt.Sprintf("%s/account/profile?key=%s", BaseURL, s.apiKey))
	if err != nil {
		return nil, err
	}

	var ret Profile
	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}
