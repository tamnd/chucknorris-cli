---
title: "Resource URIs"
description: "Use chucknorris as a database/sql-style driver so a host program can address chucknorris as chucknorris:// URIs."
weight: 20
---

`chucknorris` is a command line, but the `chucknorris` Go package is also a
small driver that makes chucknorris addressable as a resource URI. A host
program registers it the way a program registers a database driver with
`database/sql`, then dereferences `chucknorris://` URIs without knowing
anything about how chucknorris is fetched.

The host that does this today is [ant](https://github.com/tamnd/ant), a single
binary that puts one URI namespace over a family of site tools. The examples
below use `ant`; any program that links the package gets the same behaviour.

## Mounting the driver

A host enables the driver with one blank import, exactly like `import _
"github.com/lib/pq"`:

```go
import _ "github.com/tamnd/chucknorris-cli/chucknorris"
```

The package's `init` registers a domain with the scheme `chucknorris` for the
host `chucknorris.com`. The standalone `chucknorris` binary does not change.

## Addressing records

A URI is `scheme://authority/id`. The scaffold ships one type:

| URI                              | What it is                              |
| -------------------------------- | --------------------------------------- |
| `chucknorris://page/<path>`    | a page, keyed by its path on chucknorris.com |

```bash
ant get chucknorris://page/<path>    # the page record
ant cat chucknorris://page/<path>    # just the body text
ant url chucknorris://page/<path>    # the live https URL
ant resolve https://chucknorris.com/<path> # a pasted link, back to its URI
```

As you add resolver operations in `chucknorris/domain.go`, each new `URIType`
becomes another addressable authority here, with no extra wiring. See
[add a command](/guides/adding-a-command/).

## Walking the graph

`ls` lists the members of a collection, and every member is itself an
addressable URI, so a host can follow the graph and write it to disk:

```bash
ant ls     chucknorris://page/<path>             # the pages this one links to
ant export chucknorris://page/<path> --follow 1 --to ./data
```

The example `links` op emits page stubs, so each listed member is a
`chucknorris://page/` URI in its own right. When you model edges between your
real records with `kit:"link"` tags, `ant export --follow` and `ant graph` walk
those edges too, across tools when a link points at another site's scheme.

## Why this is the same code

The driver and the binary share one definition per operation. A resolver op
answers both `chucknorris page` on the command line and `ant get
chucknorris://page/...` through a host, from the same handler and the same
client. There is no second implementation to keep in step.
