package shodan

// BaseURL ...
// This is the Shodan API base url; all API functions are built off of this base URL.
// Utilized by the shodan package to more easily implement Client methods.
const BaseURL = "https://api.shodan.io"

// Client ...
// Simple struct used to hold the user's Shodan API key.
// Functions of the Shodan API are implemented as Methods to this object.
type Client struct {
	apiKey string
}

// New ...
// Returns a pointer to a new Client object.
func New(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}
