package chucknorris

import (
	"context"
	"time"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/any-cli/kit/errs"
)

// domain.go exposes chucknorris as a kit Domain driver.
//
// A multi-domain host (ant) enables it with a single blank import:
//
//	import _ "github.com/tamnd/chucknorris-cli/chucknorris"
//
// The same Domain also builds the standalone chucknorris binary.
func init() { kit.Register(Domain{}) }

// Domain is the chucknorris driver.
type Domain struct{}

// Info describes the scheme, the hostnames a pasted link is matched against,
// and the identity reused for the binary's help and version.
func (Domain) Info() kit.DomainInfo {
	return kit.DomainInfo{
		Scheme: "chucknorris",
		Hosts:  []string{Host},
		Identity: kit.Identity{
			Binary: "chucknorris",
			Short:  "Chuck Norris jokes from api.chucknorris.io",
			Long: `chucknorris fetches jokes from the public Chuck Norris joke API
at api.chucknorris.io. No API key or authentication required.`,
			Site: Host,
			Repo: "https://github.com/tamnd/chucknorris-cli",
		},
	}
}

// Register installs the client factory and every operation onto app.
func (Domain) Register(app *kit.App) {
	app.SetClient(newClient)

	// joke: get a random Chuck Norris joke
	kit.Handle(app, kit.OpMeta{
		Name:    "joke",
		Group:   "read",
		List:    false,
		Summary: "Get a random Chuck Norris joke",
	}, jokeOp)

	// search: search for jokes by query
	kit.Handle(app, kit.OpMeta{
		Name:    "search",
		Group:   "read",
		List:    true,
		Summary: "Search Chuck Norris jokes by query",
		Args:    []kit.Arg{{Name: "query", Help: "search query"}},
	}, searchOp)

	// categories: list available joke categories
	kit.Handle(app, kit.OpMeta{
		Name:    "categories",
		Group:   "read",
		List:    true,
		Summary: "List available joke categories",
	}, categoriesOp)
}

// newClient builds the client from host-resolved config.
func newClient(_ context.Context, cfg kit.Config) (any, error) {
	c := DefaultConfig()
	if cfg.UserAgent != "" {
		c.UserAgent = cfg.UserAgent
	}
	if cfg.Rate > 0 {
		c.Rate = cfg.Rate
	}
	if cfg.Retries > 0 {
		c.Retries = cfg.Retries
	}
	if cfg.Timeout > 0 {
		c.Timeout = cfg.Timeout
	}
	return NewClient(c), nil
}

// --- inputs ---

type jokeInput struct {
	Category string        `kit:"flag" help:"restrict joke to this category"`
	Delay    time.Duration `kit:"flag,inherit" help:"minimum spacing between requests"`
	Client   *Client       `kit:"inject"`
}

type searchInput struct {
	Query  string        `kit:"arg" help:"search query"`
	Limit  int           `kit:"flag,inherit" help:"max results"`
	Delay  time.Duration `kit:"flag,inherit" help:"minimum spacing between requests"`
	Client *Client       `kit:"inject"`
}

type categoriesInput struct {
	Delay  time.Duration `kit:"flag,inherit" help:"minimum spacing between requests"`
	Client *Client       `kit:"inject"`
}

// --- handlers ---

func jokeOp(ctx context.Context, in jokeInput, emit func(Joke) error) error {
	joke, err := in.Client.Random(ctx, in.Category)
	if err != nil {
		return mapErr(err)
	}
	return emit(joke)
}

func searchOp(ctx context.Context, in searchInput, emit func(Joke) error) error {
	limit := in.Limit
	if limit <= 0 {
		limit = 10
	}
	items, err := in.Client.Search(ctx, in.Query, limit)
	if err != nil {
		return mapErr(err)
	}
	for _, item := range items {
		if err := emit(item); err != nil {
			return err
		}
	}
	return nil
}

func categoriesOp(ctx context.Context, in categoriesInput, emit func(Category) error) error {
	items, err := in.Client.Categories(ctx)
	if err != nil {
		return mapErr(err)
	}
	for _, item := range items {
		if err := emit(item); err != nil {
			return err
		}
	}
	return nil
}

// --- Resolver ---

// Classify turns an input into the canonical (type, id).
func (Domain) Classify(input string) (uriType, id string, err error) {
	if input == "" {
		return "", "", errs.Usage("empty chucknorris reference")
	}
	return "joke", input, nil
}

// Locate returns the live https URL for a (type, id).
func (Domain) Locate(uriType, id string) (string, error) {
	switch uriType {
	case "joke":
		return "https://api.chucknorris.io/jokes/" + id, nil
	default:
		return "", errs.Usage("chucknorris has no resource type %q", uriType)
	}
}

// mapErr converts a library error into the kit error kind.
func mapErr(err error) error {
	return err
}
