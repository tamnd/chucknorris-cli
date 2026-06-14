package chucknorris

// Joke is one Chuck Norris joke.
type Joke struct {
	Rank  int    `json:"rank"`
	ID    string `json:"id"`
	Value string `json:"value"`
	URL   string `json:"url"`
}

// Category is one joke category returned by the categories endpoint.
type Category struct {
	Rank int    `json:"rank"`
	Name string `json:"name"`
}

// unexported: only used inside chucknorris.go for JSON decode

type rawJoke struct {
	ID      string `json:"id"`
	Value   string `json:"value"`
	IconURL string `json:"icon_url"`
	URL     string `json:"url"`
}

type searchResponse struct {
	Total  int       `json:"total"`
	Result []rawJoke `json:"result"`
}
