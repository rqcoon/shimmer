# Shimmer

Lightweight file-based CMS written in Go. Pages are written in markdown and rendered on server-side.

## Usage

Shimmer can be used as either a library or a standalone binary.

**building standalone**

```
go build -o shimmer .
./shimmer
```

**flags**

`-dir`      | Markdown content directory | default: `./page`

`-port`     | Port to listen on          | default: `8080`

`-template` | Templates directory        | default: `./templates`

**http endpoints**

Compiled pages can be requested with: `GET /page/{slug}`

Named template fragments can be requested with `frag` as a query: `GET /page/{slug}?frag=content`

The page index can be found at: `GET /index`

And a basic health endpoint exists at `/health`

**library**

The standalone binary is only a thin wrapper around the Shimmer library. A Go application can embed Shimmer directly and use its renderer, cache, templates, and HTTP handler alongside its own application routes.


## Features
- Markdown compiler
- YAML frontmatter
- Page caching
- File watcher
- Backend rendering (partial hydration for better frontend integration is WIP!)
- HTML template injection
