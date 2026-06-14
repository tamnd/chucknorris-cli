package chucknorris_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tamnd/chucknorris-cli/chucknorris"
)

const fakeJokeJSON = `{"categories":[],"created_at":"2020-01-05","icon_url":"https://assets.chucknorris.host/img/avatar/chuck-norris.png","id":"abc123","updated_at":"2020-01-05","url":"https://api.chucknorris.io/jokes/abc123","value":"Chuck Norris counted to infinity. Twice."}`

const fakeSearchJSON = `{"total":2,"result":[{"categories":[],"id":"abc123","url":"https://api.chucknorris.io/jokes/abc123","value":"Chuck Norris counted to infinity. Twice."},{"categories":[],"id":"def456","url":"https://api.chucknorris.io/jokes/def456","value":"Chuck Norris can divide by zero."}]}`

const fakeCategoriesJSON = `["animal","career","dev"]`

func newTestClient(ts *httptest.Server) *chucknorris.Client {
	cfg := chucknorris.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	return chucknorris.NewClient(cfg)
}

func TestRandomSendsUserAgent(t *testing.T) {
	var gotUA string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		_, _ = fmt.Fprint(w, fakeJokeJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.Random(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if gotUA == "" {
		t.Error("User-Agent not sent")
	}
}

func TestRandomParsesJoke(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, fakeJokeJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	joke, err := c.Random(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if joke.ID != "abc123" {
		t.Errorf("ID = %q, want abc123", joke.ID)
	}
	if joke.Value != "Chuck Norris counted to infinity. Twice." {
		t.Errorf("Value = %q, unexpected", joke.Value)
	}
	if joke.URL != "https://api.chucknorris.io/jokes/abc123" {
		t.Errorf("URL = %q, unexpected", joke.URL)
	}
}

func TestRandomWithCategory(t *testing.T) {
	var gotURL string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotURL = r.URL.String()
		_, _ = fmt.Fprint(w, fakeJokeJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.Random(context.Background(), "science")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(gotURL, "category=science") {
		t.Errorf("URL = %q, want to contain category=science", gotURL)
	}
}

func TestSearchParsesItems(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, fakeSearchJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	items, err := c.Search(context.Background(), "infinity", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}
	if items[0].ID != "abc123" {
		t.Errorf("items[0].ID = %q, want abc123", items[0].ID)
	}
	if items[1].ID != "def456" {
		t.Errorf("items[1].ID = %q, want def456", items[1].ID)
	}
}

func TestCategoriesParsesItems(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, fakeCategoriesJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	items, err := c.Categories(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 {
		t.Fatalf("len(items) = %d, want 3", len(items))
	}
	if items[0].Name != "animal" {
		t.Errorf("items[0].Name = %q, want animal", items[0].Name)
	}
	if items[2].Name != "dev" {
		t.Errorf("items[2].Name = %q, want dev", items[2].Name)
	}
}
